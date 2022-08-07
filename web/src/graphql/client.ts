import {
  GRAPHQL_SERVER_URI,
  GRAPHQL_WS_SERVER_URI,
} from "../config/enviroment";

import { ApolloClient, HttpLink, InMemoryCache, split } from "@apollo/client";
import { GraphQLWsLink } from "@apollo/client/link/subscriptions";
import { getMainDefinition } from "@apollo/client/utilities";
import { createClient } from "graphql-ws";

const httpLink = new HttpLink({
  uri: GRAPHQL_SERVER_URI,
});

const wsLink = new GraphQLWsLink(
  createClient({
    url: GRAPHQL_WS_SERVER_URI,
  }),
);

// The split function takes three parameters:
//
// * A function that's called for each operation to execute
// * The Link to use for an operation if the function returns a "truthy" value
// * The Link to use for an operation if the function returns a "falsy" value
const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);

    return (
      definition.kind === `OperationDefinition` &&
      definition.operation === `subscription`
    );
  },
  wsLink,
  httpLink,
);

const client = new ApolloClient({
  link: splitLink,
  cache: new InMemoryCache(),
});

export default client;
