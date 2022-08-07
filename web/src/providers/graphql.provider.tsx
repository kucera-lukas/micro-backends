import client from "../graphql/client";

import { ApolloProvider } from "@apollo/client";

import type { PropsWithChildren } from "react";

export type GraphqlProviderProps = PropsWithChildren<Record<never, never>>;

const GraphqlProvider = ({ children }: GraphqlProviderProps): JSX.Element => {
  return <ApolloProvider client={client}>{children}</ApolloProvider>;
};

export default GraphqlProvider;
