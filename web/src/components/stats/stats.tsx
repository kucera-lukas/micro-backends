import { PROVIDER_GROUPS } from "./constants";
import ProviderGroup from "./provider.group";

import { useProviders } from "../../context/providers.context";
import AccordionLayout from "../../layouts/accordion.layout";

import { Stack } from "@mantine/core";

const Stats = (): JSX.Element => {
  const [providers] = useProviders();

  return (
    <AccordionLayout
      value="stats"
      title="Stats"
      description="Message count statistics"
    >
      <Stack>
        <ProviderGroup
          title="Chosen providers"
          badgeColor="gray"
          providers={providers}
        />
        {PROVIDER_GROUPS.map((descriptor) => (
          <ProviderGroup
            key={descriptor.title}
            {...descriptor}
          />
        ))}
      </Stack>
    </AccordionLayout>
  );
};

export default Stats;
