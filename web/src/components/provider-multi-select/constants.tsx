import { MessageProvider } from "../../graphql/generated/codegen.generated";

import { SiMongodb, SiPostgresql } from "react-icons/all";

export const data = [
  { value: MessageProvider.Mongo, label: `Mongo` },
  { value: MessageProvider.Postgres, label: `Postgres` },
];

export const icons: Record<MessageProvider, JSX.Element> = {
  [MessageProvider.Mongo]: <SiMongodb size={20} />,
  [MessageProvider.Postgres]: <SiPostgresql size={20} />,
};
