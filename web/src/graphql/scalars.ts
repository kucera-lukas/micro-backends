import schema from "./schema";

import { withScalars } from "apollo-link-scalars";

import type { FunctionsMap } from "apollo-link-scalars";

const scalarsMap: FunctionsMap = {
  Time: {
    serialize: (parsed: unknown): string | undefined =>
      parsed instanceof Date ? parsed.toString() : undefined,
    parseValue: (raw: unknown): Date | null => {
      if (typeof raw === `string`) {
        return new Date(raw);
      }

      throw new Error(`invalid time value to parse`);
    },
  },
};

const scalarsLink = withScalars({ schema, typesMap: scalarsMap });

export default scalarsLink;
