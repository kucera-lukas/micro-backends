import { gql } from "@apollo/client";
import * as Apollo from "@apollo/client";
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  Time: Date;
};

export type Message = {
  __typename?: "Message";
  created: Scalars["Time"];
  data: Scalars["String"];
  id: Scalars["String"];
  modified: Scalars["Time"];
};

export type MessageCountPayload = {
  __typename?: "MessageCountPayload";
  count: Scalars["Int"];
  providers: Array<MessageProvider>;
};

export type MessageCreatedPayload = {
  __typename?: "MessageCreatedPayload";
  message: Message;
  provider: MessageProvider;
};

export type MessagePayload = {
  __typename?: "MessagePayload";
  message: Message;
  provider: MessageProvider;
};

export enum MessageProvider {
  Mongo = "MONGO",
  Postgres = "POSTGRES",
}

export enum MessageSortField {
  Created = "CREATED",
  Data = "DATA",
  Id = "ID",
  Modified = "MODIFIED",
}

export type MessagesPayload = {
  __typename?: "MessagesPayload";
  messages: Array<Message>;
  providers: Array<MessageProvider>;
};

/** The `Mutation` type, represents all updates we can make to our data. */
export type Mutation = {
  __typename?: "Mutation";
  newMessage: NewMessagePayload;
};

/** The `Mutation` type, represents all updates we can make to our data. */
export type MutationNewMessageArgs = {
  input: NewMessageInput;
};

export type NewMessageInput = {
  data: Scalars["String"];
  providers: Array<MessageProvider>;
};

export type NewMessagePayload = {
  __typename?: "NewMessagePayload";
  providers: Array<MessageProvider>;
  status: Scalars["String"];
};

/** The `Query` type, represents all of the entry points into our object graph. */
export type Query = {
  __typename?: "Query";
  message: MessagePayload;
  messageCount: MessageCountPayload;
  messages: MessagesPayload;
};

/** The `Query` type, represents all of the entry points into our object graph. */
export type QueryMessageArgs = {
  id: Scalars["String"];
  provider: MessageProvider;
};

/** The `Query` type, represents all of the entry points into our object graph. */
export type QueryMessageCountArgs = {
  providers: Array<MessageProvider>;
};

/** The `Query` type, represents all of the entry points into our object graph. */
export type QueryMessagesArgs = {
  providers: Array<MessageProvider>;
  reverse?: Scalars["Boolean"];
  sortField?: MessageSortField;
};

/** The `Subscription` type, represents all ways to subscribe to updates to our graph. */
export type Subscription = {
  __typename?: "Subscription";
  messageCreated: MessageCreatedPayload;
};

export type MessageFragmentFragment = {
  __typename?: "Message";
  id: string;
  data: string;
  created: Date;
  modified: Date;
};

export type NewMessageMutationVariables = Exact<{
  providers: Array<MessageProvider> | MessageProvider;
  data: Scalars["String"];
}>;

export type NewMessageMutation = {
  __typename?: "Mutation";
  newMessage: {
    __typename?: "NewMessagePayload";
    status: string;
    providers: Array<MessageProvider>;
  };
};

export type MessageCountQueryVariables = Exact<{
  providers: Array<MessageProvider> | MessageProvider;
}>;

export type MessageCountQuery = {
  __typename?: "Query";
  messageCount: {
    __typename?: "MessageCountPayload";
    count: number;
    providers: Array<MessageProvider>;
  };
};

export type MessageQueryVariables = Exact<{
  id: Scalars["String"];
  provider: MessageProvider;
}>;

export type MessageQuery = {
  __typename?: "Query";
  message: {
    __typename?: "MessagePayload";
    provider: MessageProvider;
    message: {
      __typename?: "Message";
      id: string;
      data: string;
      created: Date;
      modified: Date;
    };
  };
};

export type MessagesQueryVariables = Exact<{
  providers: Array<MessageProvider> | MessageProvider;
  sortField?: MessageSortField;
  reverse?: Scalars["Boolean"];
}>;

export type MessagesQuery = {
  __typename?: "Query";
  messages: {
    __typename?: "MessagesPayload";
    providers: Array<MessageProvider>;
    messages: Array<{
      __typename?: "Message";
      id: string;
      data: string;
      created: Date;
      modified: Date;
    }>;
  };
};

export type MessageCreatedSubscriptionVariables = Exact<{
  [key: string]: never;
}>;

export type MessageCreatedSubscription = {
  __typename?: "Subscription";
  messageCreated: {
    __typename?: "MessageCreatedPayload";
    provider: MessageProvider;
    message: {
      __typename?: "Message";
      id: string;
      data: string;
      created: Date;
      modified: Date;
    };
  };
};

export const MessageFragmentFragmentDocument = gql`
  fragment MessageFragment on Message {
    id
    data
    created
    modified
  }
`;
export const NewMessageDocument = gql`
  mutation newMessage($providers: [MessageProvider!]!, $data: String!) {
    newMessage(input: { providers: $providers, data: $data }) {
      status
      providers
    }
  }
`;
export type NewMessageMutationFn = Apollo.MutationFunction<
  NewMessageMutation,
  NewMessageMutationVariables
>;

/**
 * __useNewMessageMutation__
 *
 * To run a mutation, you first call `useNewMessageMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useNewMessageMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [newMessageMutation, { data, loading, error }] = useNewMessageMutation({
 *   variables: {
 *      providers: // value for 'providers'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useNewMessageMutation(
  baseOptions?: Apollo.MutationHookOptions<
    NewMessageMutation,
    NewMessageMutationVariables
  >,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useMutation<NewMessageMutation, NewMessageMutationVariables>(
    NewMessageDocument,
    options,
  );
}
export type NewMessageMutationHookResult = ReturnType<
  typeof useNewMessageMutation
>;
export type NewMessageMutationResult =
  Apollo.MutationResult<NewMessageMutation>;
export type NewMessageMutationOptions = Apollo.BaseMutationOptions<
  NewMessageMutation,
  NewMessageMutationVariables
>;
export const MessageCountDocument = gql`
  query messageCount($providers: [MessageProvider!]!) {
    messageCount(providers: $providers) {
      count
      providers
    }
  }
`;

/**
 * __useMessageCountQuery__
 *
 * To run a query within a React component, call `useMessageCountQuery` and pass it any options that fit your needs.
 * When your component renders, `useMessageCountQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMessageCountQuery({
 *   variables: {
 *      providers: // value for 'providers'
 *   },
 * });
 */
export function useMessageCountQuery(
  baseOptions: Apollo.QueryHookOptions<
    MessageCountQuery,
    MessageCountQueryVariables
  >,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<MessageCountQuery, MessageCountQueryVariables>(
    MessageCountDocument,
    options,
  );
}
export function useMessageCountLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    MessageCountQuery,
    MessageCountQueryVariables
  >,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<MessageCountQuery, MessageCountQueryVariables>(
    MessageCountDocument,
    options,
  );
}
export type MessageCountQueryHookResult = ReturnType<
  typeof useMessageCountQuery
>;
export type MessageCountLazyQueryHookResult = ReturnType<
  typeof useMessageCountLazyQuery
>;
export type MessageCountQueryResult = Apollo.QueryResult<
  MessageCountQuery,
  MessageCountQueryVariables
>;
export const MessageDocument = gql`
  query message($id: String!, $provider: MessageProvider!) {
    message(id: $id, provider: $provider) {
      message {
        ...MessageFragment
      }
      provider
    }
  }
  ${MessageFragmentFragmentDocument}
`;

/**
 * __useMessageQuery__
 *
 * To run a query within a React component, call `useMessageQuery` and pass it any options that fit your needs.
 * When your component renders, `useMessageQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMessageQuery({
 *   variables: {
 *      id: // value for 'id'
 *      provider: // value for 'provider'
 *   },
 * });
 */
export function useMessageQuery(
  baseOptions: Apollo.QueryHookOptions<MessageQuery, MessageQueryVariables>,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<MessageQuery, MessageQueryVariables>(
    MessageDocument,
    options,
  );
}
export function useMessageLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    MessageQuery,
    MessageQueryVariables
  >,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<MessageQuery, MessageQueryVariables>(
    MessageDocument,
    options,
  );
}
export type MessageQueryHookResult = ReturnType<typeof useMessageQuery>;
export type MessageLazyQueryHookResult = ReturnType<typeof useMessageLazyQuery>;
export type MessageQueryResult = Apollo.QueryResult<
  MessageQuery,
  MessageQueryVariables
>;
export const MessagesDocument = gql`
  query messages(
    $providers: [MessageProvider!]!
    $sortField: MessageSortField! = CREATED
    $reverse: Boolean! = true
  ) {
    messages(providers: $providers, sortField: $sortField, reverse: $reverse) {
      messages {
        ...MessageFragment
      }
      providers
    }
  }
  ${MessageFragmentFragmentDocument}
`;

/**
 * __useMessagesQuery__
 *
 * To run a query within a React component, call `useMessagesQuery` and pass it any options that fit your needs.
 * When your component renders, `useMessagesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMessagesQuery({
 *   variables: {
 *      providers: // value for 'providers'
 *      sortField: // value for 'sortField'
 *      reverse: // value for 'reverse'
 *   },
 * });
 */
export function useMessagesQuery(
  baseOptions: Apollo.QueryHookOptions<MessagesQuery, MessagesQueryVariables>,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<MessagesQuery, MessagesQueryVariables>(
    MessagesDocument,
    options,
  );
}
export function useMessagesLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    MessagesQuery,
    MessagesQueryVariables
  >,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<MessagesQuery, MessagesQueryVariables>(
    MessagesDocument,
    options,
  );
}
export type MessagesQueryHookResult = ReturnType<typeof useMessagesQuery>;
export type MessagesLazyQueryHookResult = ReturnType<
  typeof useMessagesLazyQuery
>;
export type MessagesQueryResult = Apollo.QueryResult<
  MessagesQuery,
  MessagesQueryVariables
>;
export const MessageCreatedDocument = gql`
  subscription messageCreated {
    messageCreated {
      message {
        ...MessageFragment
      }
      provider
    }
  }
  ${MessageFragmentFragmentDocument}
`;

/**
 * __useMessageCreatedSubscription__
 *
 * To run a query within a React component, call `useMessageCreatedSubscription` and pass it any options that fit your needs.
 * When your component renders, `useMessageCreatedSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMessageCreatedSubscription({
 *   variables: {
 *   },
 * });
 */
export function useMessageCreatedSubscription(
  baseOptions?: Apollo.SubscriptionHookOptions<
    MessageCreatedSubscription,
    MessageCreatedSubscriptionVariables
  >,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useSubscription<
    MessageCreatedSubscription,
    MessageCreatedSubscriptionVariables
  >(MessageCreatedDocument, options);
}
export type MessageCreatedSubscriptionHookResult = ReturnType<
  typeof useMessageCreatedSubscription
>;
export type MessageCreatedSubscriptionResult =
  Apollo.SubscriptionResult<MessageCreatedSubscription>;
