import ProviderGroup from "./provider.group";

import { useProviders } from "../../context/providers.context";
import { MessageProvider } from "../../graphql/generated/codegen.generated";

import { Accordion, Stack, Text } from "@mantine/core";

const Stats = (): JSX.Element => {
  const [providers] = useProviders();

  return (
    <Accordion variant="separated">
      <Accordion.Item value="stats">
        <Accordion.Control>
          <Text size="sm">Stats</Text>
          <Text
            size="xs"
            color="dimmed"
          >
            Message count statistics
          </Text>
        </Accordion.Control>
        <Accordion.Panel>
          <Stack>
            <ProviderGroup
              title="Chosen providers"
              badgeColor="gray"
              providers={providers}
            />
            <ProviderGroup
              title="Mongo"
              badgeColor="green"
              providers={MessageProvider.Mongo}
            />
            <ProviderGroup
              title="Postgres"
              badgeColor="blue"
              providers={MessageProvider.Postgres}
            />
          </Stack>
        </Accordion.Panel>
      </Accordion.Item>
    </Accordion>
  );
};

export default Stats;
