import { useMessageCountQuery } from "../../graphql/generated/codegen.generated";
import ErrorText from "../errors/error.text";

import { Loader, Group, Text, Badge } from "@mantine/core";

import type { MessageProvider } from "../../graphql/generated/codegen.generated";
import type { DefaultMantineColor } from "@mantine/core";

export type ProviderGroupProps = {
  title: string;
  badgeColor: DefaultMantineColor;
  providers: MessageProvider | MessageProvider[];
};

const ProviderGroup = ({
  title,
  badgeColor,
  providers,
}: ProviderGroupProps): JSX.Element => {
  const { data, loading, error } = useMessageCountQuery({
    variables: { providers },
  });

  return (
    <Group position="apart">
      <Text>{title}</Text>
      {loading ? (
        <Loader />
      ) : (
        <Badge
          color={badgeColor}
          size="lg"
        >
          {data?.messageCount.count}
        </Badge>
      )}
      {!!error && <ErrorText error={error.message} />}
    </Group>
  );
};

export default ProviderGroup;
