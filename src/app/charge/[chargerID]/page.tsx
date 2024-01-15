"use client";

import QueryChargeSession from "@/app/components/query-charge-session";
import SubscribeToChargeSession from "@/app/components/subscribe-charge-session";
import client from "@/app/contexts/graphql/graphql";
import { ApolloProvider } from "@apollo/client";

// A toggle between task 1 and task 4.
var useSubscription = false;

export default () => {
  return (
    <ApolloProvider client={client}>
      <div className="font-normal text-black font-base">
        {useSubscription ? (
          <SubscribeToChargeSession />
        ) : (
          <QueryChargeSession />
        )}
      </div>
    </ApolloProvider>
  );
};
