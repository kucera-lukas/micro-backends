import { MessageProvider } from "../../graphql/generated/codegen.generated";

import type { DefaultMantineColor } from "@mantine/core";

type ProviderDescriptor = {
  title: string;
  badgeColor: DefaultMantineColor;
  providers: MessageProvider | MessageProvider[];
};

export const PROVIDER_GROUPS: ProviderDescriptor[] = [
  {
    title: `Mongo`,
    badgeColor: `green`,
    providers: MessageProvider.Mongo,
  },
  {
    title: `Postgres`,
    badgeColor: `blue`,
    providers: MessageProvider.Postgres,
  },
];
