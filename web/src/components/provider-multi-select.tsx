import { useProviders } from "../context/providers.context";
import { MessageProvider } from "../graphql/generated/codegen.generated";
import { capitalize } from "../utils/format.utils";

import { MultiSelect } from "@mantine/core";
import { useCallback } from "react";

const data = [
  { value: MessageProvider.Mongo, label: `Mongo` },
  { value: MessageProvider.Postgres, label: `Postgres` },
];

const ProviderMultiSelect = (): JSX.Element => {
  const [providers, setProviders] = useProviders();

  const onChange = useCallback(
    (value: string[]): void => {
      setProviders(
        value.map(
          (value) =>
            MessageProvider[capitalize(value) as keyof typeof MessageProvider],
        ),
      );
    },
    [setProviders],
  );

  return (
    <MultiSelect
      value={providers}
      onChange={onChange}
      data={data}
      label="Choose message providers"
      placeholder=""
      clearable
    />
  );
};

export default ProviderMultiSelect;
