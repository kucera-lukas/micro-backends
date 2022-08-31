import { useMessageCreatedSubscription } from "../../graphql/generated/codegen.generated";

import { Center, List, Loader, Title, Card } from "@mantine/core";
import { useQueue } from "@mantine/hooks";
import { useEffect } from "react";

import type { MessageCreatedPayload } from "../../graphql/generated/codegen.generated";

const MessageFeed = (): JSX.Element => {
  const { data, loading, error } = useMessageCreatedSubscription();
  const queue = useQueue<MessageCreatedPayload>({
    initialValues: [],
    limit: 5,
  });

  console.log(queue.state);

  useEffect(() => {
    if (
      data &&
      !queue.state.includes(data.messageCreated) &&
      !queue.queue.includes(data.messageCreated)
    ) {
      console.log(`adding`);
      queue.add(data.messageCreated);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data]);

  let content;

  if (error) {
    content = <>Error: {error.message}</>;
  } else if (loading) {
    content = <Loader variant="bars" />;
  } else {
    content = (
      <List>
        {queue.state.map((messagePayload, idx) => (
          <List.Item key={idx}>
            {messagePayload.provider}: {messagePayload.message.data}
          </List.Item>
        ))}
      </List>
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
