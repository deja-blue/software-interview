"use client";

import { gql, useQuery } from "@apollo/client";
import { useParams } from "next/navigation";

type ChargeSessionModel = {
  charger: {
    id: string;
    powerKwH: number;
  };
  vehicleStateOfCharge: {
    currentBatteryLevelKwH: number;
    maxBatteryLevelKwH: number;
    rangeKmPerKwH: number;
  };
};

const CHARGE_SESSION_BY_CHARGER_ID = gql`
  query ChargeSessionByChargerID($id: String!) {
    charger: Charger(id: $id) {
      id
      powerKwH
    }
    vehicleStateOfCharge: VehicleStateOfCharge(id: $id) {
      currentBatteryLevelKwH
      maxBatteryLevelKwH
      rangeKmPerKwH
    }
  }
`;

// Task 1: Query the charger and vehicle states.
const QueryChargeSession = () => {
  const { chargerID } = useParams<{ chargerID: string }>();
  const { data, loading, error } = useQuery<ChargeSessionModel>(
    CHARGE_SESSION_BY_CHARGER_ID,
    { variables: { id: chargerID } }
  );

  return (
    <div>
      {error && <p>Error: {error.message}</p>}
      {loading && !data && <div>"Loading ..." </div>}
      <pre>
        <code>{data && JSON.stringify(data.charger, null, 2)}</code>
        <code>
          {data && JSON.stringify(data.vehicleStateOfCharge, null, 2)}
        </code>
      </pre>
    </div>
  );
};

export default QueryChargeSession;
