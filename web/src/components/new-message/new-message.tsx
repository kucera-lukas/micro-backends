import { useProviders } from "../../context/providers.context";
import { useNewMessageMutation } from "../../graphql/generated/codegen.generated";

import { TextInput, Button, Group } from "@mantine/core";
import { useCallback, useState } from "react";

const NewMessage = (): JSX.Element => {
  const [providers] = useProviders();
  const [messageData, setMessageData] = useState<string>();
  const [newMessageMutation, { loading, error }] = useNewMessageMutation();

  const onCreate = useCallback(() => {
    if (messageData !== undefined) {
      void newMessageMutation({ variables: { providers, data: messageData } });
    }
  }, [messageData, newMessageMutation, providers]);

  return (
    <Group
      grow
      spacing="xs"
    >
      <TextInput
        value={messageData}
        onChange={(event) => setMessageData(event.currentTarget.value)}
        disabled={loading}
        error={error?.message}
        label="New message"
        placeholder="Data"
        withAsterisk
      />
      <Button
        disabled={loading}
        value={messageData}
        onClick={onCreate}
        style={{ alignSelf: `flex-end` }}
      >
        Create
      </Button>
    </Group>
  );
};

export default NewMessage;
