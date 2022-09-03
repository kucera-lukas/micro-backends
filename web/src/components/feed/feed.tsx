import { FEED_MAX_SIZE, ICON_MAP } from "./constants";

import { useProviders } from "../../context/providers.context";
import { useMessageCreatedSubscription } from "../../graphql/generated/codegen.generated";
import ErrorText from "../errors/error.text";

import { Center, Loader, Stack, Group, Accordion, Text } from "@mantine/core";
import { useListState } from "@mantine/hooks";
import { useEffect } from "react";

import type { MessageCreatedPayload } from "../../graphql/generated/codegen.generated";

const Feed = (): JSX.Element => {
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

      if (messages.length >= FEED_MAX_SIZE) {
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
      <Stack>
        {messages.map((messagePayload, idx) => (
          <div key={idx}>
            <Group position="apart">
              <Text size="sm">
                {messagePayload.message.created.toLocaleTimeString()} |{` `}
                {messagePayload.message.data}
              </Text>
              {ICON_MAP[messagePayload.provider]}
            </Group>
          </div>
        ))}
        {!!error && <ErrorText error={error.message} />}
      </Stack>
    );

  return (
    <Accordion
      defaultValue="feed"
      variant="separated"
    >
      <Accordion.Item value="feed">
        <Accordion.Control>
          <Text size="sm">Feed</Text>
          <Text
            size="xs"
            color="dimmed"
          >
            Latest messages from chosen providers
          </Text>
        </Accordion.Control>
        <Accordion.Panel>{content}</Accordion.Panel>
      </Accordion.Item>
    </Accordion>
  );
};

export default Feed;
