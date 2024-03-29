import { MessageProvider } from "../../graphql/generated/codegen.generated";

import { SiMongodb, SiPostgresql } from "react-icons/all";

export const DATA = [
  { value: MessageProvider.Mongo, label: `Mongo` },
  { value: MessageProvider.Postgres, label: `Postgres` },
];

export const ICON_MAP: Record<MessageProvider, JSX.Element> = {
  [MessageProvider.Mongo]: (
    <SiMongodb
      size={20}
      color="green"
    />
  ),
  [MessageProvider.Postgres]: (
    <SiPostgresql
      size={20}
      // Mantine default blue color
      color="#a5d8ff"
    />
  ),
};
