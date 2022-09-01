const scalarsMap = {
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

export default scalarsMap;
