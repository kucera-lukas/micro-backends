import { Text } from "@mantine/core";

import type { TextProps } from "@mantine/core";

export type ErrorTextProps = {
  error: string;
  textProps?: TextProps;
};

const ErrorText = ({ error, textProps }: ErrorTextProps): JSX.Element => {
  return (
    <Text
      color="red"
      size="sm"
      mt="sm"
      {...textProps}
    >
      {error.trim()}
    </Text>
  );
};

export default ErrorText;
