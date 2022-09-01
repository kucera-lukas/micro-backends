import { icons } from "./constants";

import { Box, CloseButton } from "@mantine/core";

import type { ProviderMultiSelectValueProps } from "./types";

const ProviderMultiSelectValue = ({
  value,
  label,
  onRemove,
  ...others
}: ProviderMultiSelectValueProps): JSX.Element => {
  const icon = icons[value];

  return (
    <div {...others}>
      <Box
        sx={(theme) => ({
          display: `flex`,
          cursor: `default`,
          alignItems: `center`,
          backgroundColor:
            theme.colorScheme === `dark` ? theme.colors.dark[7] : theme.white,
          border: `1px solid ${
            theme.colorScheme === `dark`
              ? theme.colors.dark[7]
              : theme.colors.gray[4]
          }`,
          paddingLeft: 10,
          borderRadius: 4,
        })}
      >
        <Box
          mr={5}
          mt={5}
        >
          {icon}
        </Box>
        <Box sx={{ lineHeight: 1, fontSize: 12 }}>{label}</Box>
        <CloseButton
          onMouseDown={onRemove}
          variant="transparent"
          size={22}
          iconSize={14}
          tabIndex={-1}
        />
      </Box>
    </div>
  );
};

export default ProviderMultiSelectValue;
