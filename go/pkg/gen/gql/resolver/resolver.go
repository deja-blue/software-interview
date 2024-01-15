package resolver

import "github.com/deja-blue/swe-interview/go/pkg/charge"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ChargeResolver *charge.Resolver
}
