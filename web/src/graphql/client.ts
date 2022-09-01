import scalarsMap from "./scalars";

import {
  GRAPHQL_SERVER_URI,
  GRAPHQL_WS_SERVER_URI,
} from "../config/enviroment";
import introspectionSchema from "../graphql/generated/introspection-schema.generated.json";

import {
  ApolloClient,
  ApolloLink,
  HttpLink,
  InMemoryCache,
  split,
} from "@apollo/client";
import { GraphQLWsLink } from "@apollo/client/link/subscriptions";
import { getMainDefinition } from "@apollo/client/utilities";
import { withScalars } from "apollo-link-scalars";
import { buildClientSchema } from "graphql";
import { createClient } from "graphql-ws";

import type { IntrospectionQuery } from "graphql";

const schema = buildClientSchema(
  introspectionSchema as unknown as IntrospectionQuery,
);

const scalarsLink = withScalars({ schema, typesMap: scalarsMap });

const httpLink = ApolloLink.from([
  scalarsLink,
  new HttpLink({
    uri: GRAPHQL_SERVER_URI,
  }),
]);

const wsLink = ApolloLink.from([
  scalarsLink,
  new GraphQLWsLink(
    createClient({
      url: GRAPHQL_WS_SERVER_URI,
    }),
  ),
]);

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
