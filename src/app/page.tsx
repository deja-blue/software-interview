"use client";

import { ApolloProvider } from "@apollo/client";
import client from "./contexts/graphql/graphql";
import PairCharger from "./components/swe-start";

function SWEInterview({}) {
  return (
    <ApolloProvider client={client}>
      <PairCharger />
    </ApolloProvider>
  );
}

export default SWEInterview;
