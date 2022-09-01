import type { MessageProvider } from "../../graphql/generated/codegen.generated";
import type { MultiSelectValueProps, SelectItemProps } from "@mantine/core";

export interface ProviderMultiSelectItemProps extends SelectItemProps {
  value?: MessageProvider;
}

export interface ProviderMultiSelectValueProps extends MultiSelectValueProps {
  value: MessageProvider;
}
