import Feed from "../components/feed";
import NewMessage from "../components/new-message";
import ProviderMultiSelect from "../components/provider-multi-select";
import Stats from "../components/stats";

import { Center, Stack, Title } from "@mantine/core";

const Home = (): JSX.Element => {
  return (
    <Center>
      <Stack>
        <Title>micro-backends</Title>
        <ProviderMultiSelect />
        <Feed />
        <Stats />
        <NewMessage />
      </Stack>
    </Center>
  );
};

export default Home;
