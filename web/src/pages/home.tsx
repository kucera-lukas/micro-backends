import ProviderMultiSelect from "../components/provider-multi-select";
import {
  MessageProvider,
  useMessageQuery,
} from "../graphql/generated/codegen.generated";

import { Center, Stack, Title } from "@mantine/core";

const Home = (): JSX.Element => {
  const { data, loading, error } = useMessageQuery({
    variables: {
      id: `1`,
      provider: MessageProvider.Postgres,
    },
  });

  let content;
  if (loading) {
    content = <div>Loading...</div>;
  } else if (error) {
    content = <div>Error {error.message}</div>;
  } else {
    content = <div>{data?.message.message.data}</div>;
  }

  return (
    <Center>
      <Stack>
        <Title>Micro Backends</Title>
        <ProviderMultiSelect />
        <div>{content}</div>
      </Stack>
    </Center>
  );
};

export default Home;
