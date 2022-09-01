import introspectionSchema from "./generated/introspection-schema.generated.json";

import { buildClientSchema } from "graphql";

import type { IntrospectionQuery } from "graphql";

const schema = buildClientSchema(
  introspectionSchema as unknown as IntrospectionQuery,
);

export default schema;
