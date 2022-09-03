import { useProviders } from "../../context/providers.context";
import { useNewMessageMutation } from "../../graphql/generated/codegen.generated";
import AccordionLayout from "../../layouts/accordion.layout";

import { TextInput, Button, Stack } from "@mantine/core";
import { useCallback, useEffect, useState } from "react";

const NewMessage = (): JSX.Element => {
  const [error, setError] = useState<string>();
  const [providers] = useProviders();
  const [messageData, setMessageData] = useState<string>();
  const [newMessageMutation, { loading, error: mutationError }] =
    useNewMessageMutation({
      refetchQueries: [`messageCount`],
    });

  const onCreate = useCallback(() => {
    if (providers.length === 0) {
      setError(`Choose at least one provider`);
    } else if (messageData === undefined || messageData?.length === 0) {
      setError(`Message data is required`);
    } else {
      // eslint-disable-next-line unicorn/no-useless-undefined
      setError(undefined);
      void newMessageMutation({ variables: { providers, data: messageData } });
    }
  }, [messageData, newMessageMutation, providers]);

  useEffect(() => {
    if (mutationError) {
      setError(mutationError.message);
    }
  }, [mutationError]);

  return (
    <AccordionLayout
      value="new-message"
      title="New message"
    >
      <Stack spacing="xs">
        <TextInput
          value={messageData}
          onChange={(event) => setMessageData(event.currentTarget.value)}
          disabled={loading}
          error={error}
          label="Data"
          description="Message will be created by each provider"
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
      </Stack>
    </AccordionLayout>
  );
};

export default NewMessage;
