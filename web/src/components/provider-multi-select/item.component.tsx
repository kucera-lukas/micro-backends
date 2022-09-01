import { icons } from "./constants";

import { Box } from "@mantine/core";
import { forwardRef } from "react";

import type { ProviderMultiSelectItemProps } from "./types";

const ProviderMultiSelectItem = forwardRef<
  HTMLDivElement,
  ProviderMultiSelectItemProps
>(({ label, value, ...others }, ref) => {
  const icon = value && icons[value];

  return (
    <div
      ref={ref}
      {...others}
    >
      <Box sx={{ display: `flex`, alignItems: `center` }}>
        <Box mr={5}>{icon}</Box>
        <div>{label}</div>
      </Box>
    </div>
  );
});

ProviderMultiSelectItem.displayName = `Item`;

export default ProviderMultiSelectItem;
