import { useProviders } from "../../context/providers.context";
import { useMessageCreatedSubscription } from "../../graphql/generated/codegen.generated";

import { Center, Loader, Title, Card, Stack } from "@mantine/core";
import { useListState } from "@mantine/hooks";
import { useEffect } from "react";

import type { MessageCreatedPayload } from "../../graphql/generated/codegen.generated";

const MessageFeed = (): JSX.Element => {
  const { data, loading, error } = useMessageCreatedSubscription();
  const [messages, messageHandlers] = useListState<MessageCreatedPayload>([]);
  const [providers] = useProviders();

  useEffect(() => {
    if (
      data &&
      // record this message only if its provider is selected
      providers.includes(data.messageCreated.provider)
    ) {
      messageHandlers.append(data.messageCreated);

      while (messages.length > 5) {
        messageHandlers.shift();
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data]);

  let content;

  if (error) {
    content = <>Error: {error.message}</>;
  } else if (loading || messages.length === 0) {
    content = <Loader variant="bars" />;
  } else {
    content = (
      <Stack justify="flex-end">
        {messages.map((messagePayload, idx) => (
          <div key={idx}>
            <>
              [{messagePayload.provider}]{` `}
              {/*eslint-disable-next-line @typescript-eslint/ban-ts-comment*/}
              {/*@ts-ignore*/}
              {new Date(messagePayload.message.created).toLocaleTimeString()}:
              {` `}
              {messagePayload.message.data}
            </>
          </div>
        ))}
      </Stack>
    );
  }

  return (
    <Card>
      <Card.Section
        withBorder
        inheritPadding
        py="xs"
      >
        <Title order={5}>Message Feed</Title>
      </Card.Section>
      <Card.Section
        withBorder
        inheritPadding
        py="xs"
      >
        <Center>{content}</Center>
      </Card.Section>
    </Card>
  );
};

export default MessageFeed;
