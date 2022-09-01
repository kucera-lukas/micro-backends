import MessageFeed from "../components/message-feed";
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

  let content: JSX.Element;
  if (loading) {
    content = <>Loading...</>;
  } else if (error) {
    content = <>Error {error.message}</>;
  } else {
    content = <>{data?.message.message.data}</>;
  }

  return (
    <Center>
      <Stack>
        <Title>micro-backends</Title>
        <ProviderMultiSelect />
        <MessageFeed />
        <div>{content}</div>
      </Stack>
    </Center>
  );
};

export default Home;
