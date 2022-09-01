import { FEED_MAX_SIZE } from "./constants";

import { useProviders } from "../../context/providers.context";
import { useMessageCreatedSubscription } from "../../graphql/generated/codegen.generated";

import { Center, Loader, Stack, Text, Group, Accordion } from "@mantine/core";
import { useListState } from "@mantine/hooks";
import { useEffect } from "react";

import type { MessageCreatedPayload } from "../../graphql/generated/codegen.generated";

const MessageFeed = (): JSX.Element => {
  const { data, loading, error } = useMessageCreatedSubscription();
  // we don't need to use more efficient data structure as we only operate with
  // 'FEED_MAX_SIZE' number of elements
  const [messages, messageHandlers] = useListState<MessageCreatedPayload>([]);
  const [providers] = useProviders();

  useEffect(() => {
    if (
      data &&
      // record this message only if its provider is selected
      providers.includes(data.messageCreated.provider)
    ) {
      messageHandlers.append(data.messageCreated);

      if (messages.length > FEED_MAX_SIZE) {
        messageHandlers.shift();
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data]);

  const content =
    loading || messages.length === 0 ? (
      <Center>
        <Loader variant="bars" />
      </Center>
    ) : (
      <Group position="left">
        <Stack justify="flex-end">
          {messages.map((messagePayload, idx) => (
            <div key={idx}>
              <>
                [{messagePayload.provider} - {messagePayload.message.id}]{` `}
                {messagePayload.message.created.toLocaleTimeString()}:{` `}
                {messagePayload.message.data}
              </>
            </div>
          ))}
        </Stack>
        {!!error && (
          <Text
            color="red"
            size="sm"
            mt="sm"
          >
            Error: {error.message}
          </Text>
        )}
      </Group>
    );

  return (
    <Accordion
      defaultValue="feed"
      variant="separated"
    >
      <Accordion.Item value="feed">
        <Accordion.Control>Message Feed</Accordion.Control>
        <Accordion.Panel>{content}</Accordion.Panel>
      </Accordion.Item>
    </Accordion>
  );
};

export default MessageFeed;
