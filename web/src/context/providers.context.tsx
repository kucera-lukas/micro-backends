import { createContext, useContext, useState } from "react";

import type { MessageProvider } from "../graphql/generated/codegen.generated";
import type { PropsWithChildren, Dispatch, SetStateAction } from "react";

export type MessageProviderPayload = readonly [
  MessageProvider[],
  Dispatch<SetStateAction<MessageProvider[]>>,
];

export const MessageProvidersContext = createContext<
  MessageProviderPayload | undefined
>(undefined);

export type MessageProvidersProviderProps = PropsWithChildren<
  Record<never, never>
>;

export const MessageProvidersProvider = ({
  children,
}: MessageProvidersProviderProps): JSX.Element => {
  const [providers, setProviders] = useState<MessageProvider[]>([]);

  return (
    <MessageProvidersContext.Provider value={[providers, setProviders]}>
      {children}
    </MessageProvidersContext.Provider>
  );
};

export const useProviders = (): MessageProviderPayload => {
  const context = useContext(MessageProvidersContext);

  if (!context) {
    throw new Error(
      `useProviders must be used within a MessageProvidersProvider`,
    );
  }

  return context;
};
