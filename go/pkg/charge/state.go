package charge

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/deja-blue/swe-interview/go/pkg/gen/gql/model"
	"github.com/google/uuid"
)

var (
	ErrChargerNotOccupied = errors.New("charger is not occupied")
)

type ChargerID string

type ChargerState string

const (
	Available ChargerState = "available"
	Occupied  ChargerState = "occupied"
)

//// Data models for object representation within the codebsae.
// Unexported types below (ex. chargerState) are used to represent the database version of the object.
// External API boundary models are represented in the gqlgen models, in the go/pkg/gen/gql/model package.

// Charger represents a charger.
// This represents the queried object from a database, which would be constant for the lifecyle of a request or other operation.
type Charger struct {
	ID                   ChargerID
	MaxKwHDraw           float64
	vehicleStateOfCharge *VehicleStateOfCharge
}

// VehicleStateOfCharge represents the database version of the vehicle state of charge.
type VehicleStateOfCharge struct {
	CurrentBatteryLevelKwH float64
	MaxBatteryLevelKwH     float64
	RangeKmPerKwH          float64
}

func (c *Charger) Status() model.ChargerStatus {
	if c.vehicleStateOfCharge != nil {
		return model.ChargerStatusPluggedIn
	}
	return model.ChargerStatusAvailable
}

func (c *Charger) ToChargeModel() *model.Charger {
	chargeModel := &model.Charger{
		ID:       string(c.ID),
		PowerKwH: int(c.MaxKwHDraw),
	}

	return chargeModel
}

func (c *Charger) ToVehicleStateOfChargeModel() (*model.VehicleStateOfCharge, error) {
	if c.vehicleStateOfCharge == nil {
		return nil, ErrChargerNotOccupied
	}

	maxBatteryLevel := c.vehicleStateOfCharge.MaxBatteryLevelKwH
	currentBatteryLevel := c.vehicleStateOfCharge.CurrentBatteryLevelKwH
	rangeKmPerKwH := c.vehicleStateOfCharge.RangeKmPerKwH
	return &model.VehicleStateOfCharge{
		CurrentBatteryLevelKwH: &currentBatteryLevel,
		MaxBatteryLevelKwH:     &maxBatteryLevel,
		RangeKmPerKwH:          &rangeKmPerKwH,
	}, nil
}

// chargerState represents the state of a charger.
// In a real life application, this would be stored in a database.
// For this exercise, we will store it in memory.
// The updater is acting as a simulator to represent updating charger state from external events.
type chargerState struct {
	ID                   ChargerID
	MaxKwHDraw           float64
	vehicleStateOfCharge *VehicleStateOfCharge
	m                    sync.RWMutex
	chargerStateUpdater  *chargerStateUpdater
	subscribers          map[uuid.UUID]chan struct{}
}

type chargerStateUpdater struct {
	// This is a simulator to represent updating charger state from external events.
	// This is used to enable the mock to change state BetweenPluggedIn and Available.
	// If this is set, the charger will be in available.
	// On update, we will reset the vehicle state of charge to to this value and clear the field,
	// thus enabling a simulator for a vehicle being plugged/unplugged from a charger
	vehicleStateOfChargeStore *VehicleStateOfCharge

	updaterFrequency time.Duration
}

func (c *chargerState) Export() *Charger {
	c.m.RLock()
	defer c.m.RUnlock()

	return &Charger{
		ID:                   c.ID,
		MaxKwHDraw:           c.MaxKwHDraw,
		vehicleStateOfCharge: c.vehicleStateOfCharge,
	}
}

// subscribToStateUpdates will return a channel that will receive updates when the charger state changes.
// The channel (`subscriber`) must be removed when the context is canceled, (otherwise one will not be able to resubscribe).
// Hint: this operation must not block
// The channel will receive an update when the charger state changes,
// (so look at where the update happens and send the message to the channel there)
//
// For the demo, we will only allow one subscriber.
// Rememmber to handle concurrency and ensure safe updates to Charger properties.
func (c *chargerState) subscribToStateUpdates(ctx context.Context) (<-chan struct{}, error) {
	ch := make(chan struct{}, 1)
	id := uuid.New()
	c.m.Lock()
	defer c.m.Unlock()

	c.subscribers[id] = ch
	go func() {
		<-ctx.Done()
		c.m.Lock()
		defer c.m.Unlock()
		c.subscribers[id] = nil
	}()
	return ch, nil
}

func (c *chargerState) receiveChargerStateUpdates() {
	// We will only receive updates if the mock starts with vehicle inofrmation to toggle back and forth.
	if c.chargerStateUpdater == nil || c.vehicleStateOfCharge == nil {
		return
	}

	go func() {
		for {
			select {
			case <-time.After(c.chargerStateUpdater.updaterFrequency):
				c.m.RLock()
				if c.chargerStateUpdater.vehicleStateOfChargeStore != nil {
					c.vehicleStateOfCharge = c.chargerStateUpdater.vehicleStateOfChargeStore
					c.chargerStateUpdater.vehicleStateOfChargeStore = nil
				} else {
					c.chargerStateUpdater.vehicleStateOfChargeStore = c.vehicleStateOfCharge
					c.vehicleStateOfCharge = nil
				}
				for _, sub := range c.subscribers {
					sub := sub
					go func() {
						select {
						case sub <- struct{}{}:
						default:
							// Subscriber channel is not ready to receive the update
						}
					}()
				}
				c.m.RUnlock()
			}
		}
	}()
}

type MockChargerOption func(*chargerState)

func MockCharger(id string, availablePowerKwH float64, opts ...MockChargerOption) *chargerState {
	c := &chargerState{
		ID:          ChargerID(id),
		MaxKwHDraw:  availablePowerKwH,
		subscribers: make(map[uuid.UUID]chan struct{}),
	}
	for _, o := range opts {
		o(c)
	}
	c.receiveChargerStateUpdates()
	return c
}

func WithStateOfCharge(stateOfCharge *VehicleStateOfCharge) MockChargerOption {
	return func(c *chargerState) {
		c.vehicleStateOfCharge = stateOfCharge
	}
}

func WithDynamicStatus(updateFrequency time.Duration) MockChargerOption {
	return func(c *chargerState) {
		c.chargerStateUpdater = &chargerStateUpdater{
			updaterFrequency: updateFrequency,
		}
	}
}

type Resolver struct {
	chargers map[ChargerID]*chargerState
}

func (r *Resolver) SubscribToStateUpdates(ctx context.Context, id ChargerID) (<-chan struct{}, error) {
	c, ok := r.chargers[id]
	if !ok {
		return nil, fmt.Errorf("charger not found for id %s", id)
	}
	return c.subscribToStateUpdates(ctx)
}

func (r *Resolver) Charger(id ChargerID) (*Charger, error) {
	c, ok := r.chargers[id]
	if !ok {
		return nil, fmt.Errorf("charger not found for id %s", id)
	}
	return c.Export(), nil
}

func NewResolver(chargers ...*chargerState) *Resolver {
	m := make(map[ChargerID]*chargerState)
	for _, c := range chargers {
		m[c.ID] = c
	}
	return &Resolver{
		chargers: m,
	}
}
