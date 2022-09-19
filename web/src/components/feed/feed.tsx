import { FEED_MAX_SIZE, ICON_MAP } from "./constants";

import { useProviders } from "../../context/providers.context";
import client from "../../graphql/client";
import { useMessageCreatedSubscription } from "../../graphql/generated/codegen.generated";
import AccordionLayout from "../../layouts/accordion.layout";
import ErrorText from "../errors/error.text";

import { Center, Loader, Stack, Group, Text } from "@mantine/core";
import { useListState } from "@mantine/hooks";
import { useCallback } from "react";

import type {
  MessageCreatedPayload,
  MessageCreatedSubscription,
} from "../../graphql/generated/codegen.generated";
import type { OnSubscriptionDataOptions } from "@apollo/client";

const Feed = (): JSX.Element => {
  // we don't need to use more efficient data structure as we only operate with
  // 'FEED_MAX_SIZE' number of elements
  const [messages, messageHandlers] = useListState<MessageCreatedPayload>([]);
  const [providers] = useProviders();
  const { loading, error } = useMessageCreatedSubscription({
    onSubscriptionData: useCallback(
      (data: OnSubscriptionDataOptions<MessageCreatedSubscription>) => {
        void client.refetchQueries({ include: [`messageCount`] });

        const messagePayload = data.subscriptionData.data?.messageCreated;

        if (
          messagePayload &&
          // record this message only if its provider is selected
          providers.includes(messagePayload.provider)
        ) {
          messageHandlers.append(messagePayload);

          if (messages.length >= FEED_MAX_SIZE) {
            messageHandlers.shift();
          }
        }
      },
      [messageHandlers, messages.length, providers],
    ),
  });

  return (
    <AccordionLayout
      value="feed"
      title="Feed"
      description="Latest messages from chosen providers"
    >
      {loading || messages.length === 0 ? (
        <Center>
          <Loader variant="bars" />
        </Center>
      ) : (
        <Stack>
          {messages.map((messagePayload, idx) => (
            <Group
              key={idx}
              position="apart"
            >
              <Text size="sm">
                {messagePayload.message.created.toLocaleTimeString()} |{` `}
                {messagePayload.message.data}
              </Text>
              {ICON_MAP[messagePayload.provider]}
            </Group>
          ))}
          {!!error && <ErrorText error={error.message} />}
        </Stack>
      )}
    </AccordionLayout>
  );
};

export default Feed;
