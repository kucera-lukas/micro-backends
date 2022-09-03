import { Accordion, Text } from "@mantine/core";

import type { PropsWithChildren } from "react";

export type AccordionLayoutProps = PropsWithChildren<{
  value: string;
  title: string;
  description?: string;
}>;

const AccordionLayout = ({
  value,
  title,
  description,
  children,
}: AccordionLayoutProps): JSX.Element => {
  return (
    <Accordion
      defaultValue={value}
      variant="separated"
    >
      <Accordion.Item value={value}>
        <Accordion.Control>
          <Text size="sm">{title}</Text>
          {description && (
            <Text
              size="xs"
              color="dimmed"
            >
              {description}
            </Text>
          )}
        </Accordion.Control>
        <Accordion.Panel>{children}</Accordion.Panel>
      </Accordion.Item>
    </Accordion>
  );
};

export default AccordionLayout;
