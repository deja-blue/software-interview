"use client";
import { useParams } from "next/navigation";
import { gql, TypedDocumentNode, useSubscription } from "@apollo/client";

// Task 4: Subscribe to charger status updates.
const SUBSCRIBE_CHARGER_UPDATES = gql`
  subscription SubscribeToChargerStateChanges($id: String!) {
    ChargerState(id: $id) {
      chargerStatus
    }
  }
` as TypedDocumentNode<ChargeSubscriptionModel>;

type ChargeSubscriptionModel = {
  ChargerState: {
    // To update with the appropriate terms for the exercise
    chargerStatus: string;
  };
};

const SubscribeToChargeSession = () => {
  const { chargerID } = useParams<{ chargerID: string }>();
  const { data, loading, error } = useSubscription(SUBSCRIBE_CHARGER_UPDATES, {
    variables: { id: chargerID },
  });

  return (
    <div>
      {error && <p>Error: {error.message}</p>}
      {loading && !data && <div>"Loading ..." </div>}
      {data && data.ChargerState && (
        <div>ChargerStatus: {data.ChargerState.chargerStatus} </div>
      )}
    </div>
  );
};

export default SubscribeToChargeSession;
