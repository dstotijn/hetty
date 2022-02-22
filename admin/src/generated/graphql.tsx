import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  Regexp: any;
  Time: any;
  URL: any;
};

export type ClearHttpRequestLogResult = {
  __typename?: 'ClearHTTPRequestLogResult';
  success: Scalars['Boolean'];
};

export type CloseProjectResult = {
  __typename?: 'CloseProjectResult';
  success: Scalars['Boolean'];
};

export type DeleteProjectResult = {
  __typename?: 'DeleteProjectResult';
  success: Scalars['Boolean'];
};

export type DeleteSenderRequestsResult = {
  __typename?: 'DeleteSenderRequestsResult';
  success: Scalars['Boolean'];
};

export type HttpHeader = {
  __typename?: 'HttpHeader';
  key: Scalars['String'];
  value: Scalars['String'];
};

export type HttpHeaderInput = {
  key: Scalars['String'];
  value: Scalars['String'];
};

export enum HttpMethod {
  Connect = 'CONNECT',
  Delete = 'DELETE',
  Get = 'GET',
  Head = 'HEAD',
  Options = 'OPTIONS',
  Patch = 'PATCH',
  Post = 'POST',
  Put = 'PUT',
  Trace = 'TRACE'
}

export enum HttpProtocol {
  Http1 = 'HTTP1',
  Http2 = 'HTTP2'
}

export type HttpRequestLog = {
  __typename?: 'HttpRequestLog';
  body?: Maybe<Scalars['String']>;
  headers: Array<HttpHeader>;
  id: Scalars['ID'];
  method: HttpMethod;
  proto: Scalars['String'];
  response?: Maybe<HttpResponseLog>;
  timestamp: Scalars['Time'];
  url: Scalars['String'];
};

export type HttpRequestLogFilter = {
  __typename?: 'HttpRequestLogFilter';
  onlyInScope: Scalars['Boolean'];
  searchExpression?: Maybe<Scalars['String']>;
};

export type HttpRequestLogFilterInput = {
  onlyInScope?: InputMaybe<Scalars['Boolean']>;
  searchExpression?: InputMaybe<Scalars['String']>;
};

export type HttpResponseLog = {
  __typename?: 'HttpResponseLog';
  body?: Maybe<Scalars['String']>;
  headers: Array<HttpHeader>;
  /** Will be the same ID as its related request ID. */
  id: Scalars['ID'];
  proto: HttpProtocol;
  statusCode: Scalars['Int'];
  statusReason: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  clearHTTPRequestLog: ClearHttpRequestLogResult;
  closeProject: CloseProjectResult;
  createOrUpdateSenderRequest: SenderRequest;
  createProject?: Maybe<Project>;
  createSenderRequestFromHttpRequestLog: SenderRequest;
  deleteProject: DeleteProjectResult;
  deleteSenderRequests: DeleteSenderRequestsResult;
  openProject?: Maybe<Project>;
  sendRequest: SenderRequest;
  setHttpRequestLogFilter?: Maybe<HttpRequestLogFilter>;
  setScope: Array<ScopeRule>;
  setSenderRequestFilter?: Maybe<SenderRequestFilter>;
};


export type MutationCreateOrUpdateSenderRequestArgs = {
  request: SenderRequestInput;
};


export type MutationCreateProjectArgs = {
  name: Scalars['String'];
};


export type MutationCreateSenderRequestFromHttpRequestLogArgs = {
  id: Scalars['ID'];
};


export type MutationDeleteProjectArgs = {
  id: Scalars['ID'];
};


export type MutationOpenProjectArgs = {
  id: Scalars['ID'];
};


export type MutationSendRequestArgs = {
  id: Scalars['ID'];
};


export type MutationSetHttpRequestLogFilterArgs = {
  filter?: InputMaybe<HttpRequestLogFilterInput>;
};


export type MutationSetScopeArgs = {
  scope: Array<ScopeRuleInput>;
};


export type MutationSetSenderRequestFilterArgs = {
  filter?: InputMaybe<SenderRequestFilterInput>;
};

export type Project = {
  __typename?: 'Project';
  id: Scalars['ID'];
  isActive: Scalars['Boolean'];
  name: Scalars['String'];
};

export type Query = {
  __typename?: 'Query';
  activeProject?: Maybe<Project>;
  httpRequestLog?: Maybe<HttpRequestLog>;
  httpRequestLogFilter?: Maybe<HttpRequestLogFilter>;
  httpRequestLogs: Array<HttpRequestLog>;
  projects: Array<Project>;
  scope: Array<ScopeRule>;
  senderRequest?: Maybe<SenderRequest>;
  senderRequests: Array<SenderRequest>;
};


export type QueryHttpRequestLogArgs = {
  id: Scalars['ID'];
};


export type QuerySenderRequestArgs = {
  id: Scalars['ID'];
};

export type ScopeHeader = {
  __typename?: 'ScopeHeader';
  key?: Maybe<Scalars['Regexp']>;
  value?: Maybe<Scalars['Regexp']>;
};

export type ScopeHeaderInput = {
  key?: InputMaybe<Scalars['Regexp']>;
  value?: InputMaybe<Scalars['Regexp']>;
};

export type ScopeRule = {
  __typename?: 'ScopeRule';
  body?: Maybe<Scalars['Regexp']>;
  header?: Maybe<ScopeHeader>;
  url?: Maybe<Scalars['Regexp']>;
};

export type ScopeRuleInput = {
  body?: InputMaybe<Scalars['Regexp']>;
  header?: InputMaybe<ScopeHeaderInput>;
  url?: InputMaybe<Scalars['Regexp']>;
};

export type SenderRequest = {
  __typename?: 'SenderRequest';
  body?: Maybe<Scalars['String']>;
  headers?: Maybe<Array<HttpHeader>>;
  id: Scalars['ID'];
  method: HttpMethod;
  proto: HttpProtocol;
  response?: Maybe<HttpResponseLog>;
  sourceRequestLogID?: Maybe<Scalars['ID']>;
  timestamp: Scalars['Time'];
  url: Scalars['URL'];
};

export type SenderRequestFilter = {
  __typename?: 'SenderRequestFilter';
  onlyInScope: Scalars['Boolean'];
  searchExpression?: Maybe<Scalars['String']>;
};

export type SenderRequestFilterInput = {
  onlyInScope?: InputMaybe<Scalars['Boolean']>;
  searchExpression?: InputMaybe<Scalars['String']>;
};

export type SenderRequestInput = {
  body?: InputMaybe<Scalars['String']>;
  headers?: InputMaybe<Array<HttpHeaderInput>>;
  id?: InputMaybe<Scalars['ID']>;
  method?: InputMaybe<HttpMethod>;
  proto?: InputMaybe<HttpProtocol>;
  url: Scalars['URL'];
};

export type CreateOrUpdateSenderRequestMutationVariables = Exact<{
  request: SenderRequestInput;
}>;


export type CreateOrUpdateSenderRequestMutation = { __typename?: 'Mutation', createOrUpdateSenderRequest: { __typename?: 'SenderRequest', id: string } };

export type CreateSenderRequestFromHttpRequestLogMutationVariables = Exact<{
  id: Scalars['ID'];
}>;


export type CreateSenderRequestFromHttpRequestLogMutation = { __typename?: 'Mutation', createSenderRequestFromHttpRequestLog: { __typename?: 'SenderRequest', id: string } };

export type SendRequestMutationVariables = Exact<{
  id: Scalars['ID'];
}>;


export type SendRequestMutation = { __typename?: 'Mutation', sendRequest: { __typename?: 'SenderRequest', id: string } };

export type GetSenderRequestQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetSenderRequestQuery = { __typename?: 'Query', senderRequest?: { __typename?: 'SenderRequest', id: string, sourceRequestLogID?: string | null, url: any, method: HttpMethod, proto: HttpProtocol, body?: string | null, timestamp: any, headers?: Array<{ __typename?: 'HttpHeader', key: string, value: string }> | null, response?: { __typename?: 'HttpResponseLog', id: string, proto: HttpProtocol, statusCode: number, statusReason: string, body?: string | null, headers: Array<{ __typename?: 'HttpHeader', key: string, value: string }> } | null } | null };

export type GetSenderRequestsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetSenderRequestsQuery = { __typename?: 'Query', senderRequests: Array<{ __typename?: 'SenderRequest', id: string, url: any, method: HttpMethod, response?: { __typename?: 'HttpResponseLog', id: string, statusCode: number, statusReason: string } | null }> };


export const CreateOrUpdateSenderRequestDocument = gql`
    mutation CreateOrUpdateSenderRequest($request: SenderRequestInput!) {
  createOrUpdateSenderRequest(request: $request) {
    id
  }
}
    `;
export type CreateOrUpdateSenderRequestMutationFn = Apollo.MutationFunction<CreateOrUpdateSenderRequestMutation, CreateOrUpdateSenderRequestMutationVariables>;

/**
 * __useCreateOrUpdateSenderRequestMutation__
 *
 * To run a mutation, you first call `useCreateOrUpdateSenderRequestMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateOrUpdateSenderRequestMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createOrUpdateSenderRequestMutation, { data, loading, error }] = useCreateOrUpdateSenderRequestMutation({
 *   variables: {
 *      request: // value for 'request'
 *   },
 * });
 */
export function useCreateOrUpdateSenderRequestMutation(baseOptions?: Apollo.MutationHookOptions<CreateOrUpdateSenderRequestMutation, CreateOrUpdateSenderRequestMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateOrUpdateSenderRequestMutation, CreateOrUpdateSenderRequestMutationVariables>(CreateOrUpdateSenderRequestDocument, options);
      }
export type CreateOrUpdateSenderRequestMutationHookResult = ReturnType<typeof useCreateOrUpdateSenderRequestMutation>;
export type CreateOrUpdateSenderRequestMutationResult = Apollo.MutationResult<CreateOrUpdateSenderRequestMutation>;
export type CreateOrUpdateSenderRequestMutationOptions = Apollo.BaseMutationOptions<CreateOrUpdateSenderRequestMutation, CreateOrUpdateSenderRequestMutationVariables>;
export const CreateSenderRequestFromHttpRequestLogDocument = gql`
    mutation CreateSenderRequestFromHttpRequestLog($id: ID!) {
  createSenderRequestFromHttpRequestLog(id: $id) {
    id
  }
}
    `;
export type CreateSenderRequestFromHttpRequestLogMutationFn = Apollo.MutationFunction<CreateSenderRequestFromHttpRequestLogMutation, CreateSenderRequestFromHttpRequestLogMutationVariables>;

/**
 * __useCreateSenderRequestFromHttpRequestLogMutation__
 *
 * To run a mutation, you first call `useCreateSenderRequestFromHttpRequestLogMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateSenderRequestFromHttpRequestLogMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createSenderRequestFromHttpRequestLogMutation, { data, loading, error }] = useCreateSenderRequestFromHttpRequestLogMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useCreateSenderRequestFromHttpRequestLogMutation(baseOptions?: Apollo.MutationHookOptions<CreateSenderRequestFromHttpRequestLogMutation, CreateSenderRequestFromHttpRequestLogMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateSenderRequestFromHttpRequestLogMutation, CreateSenderRequestFromHttpRequestLogMutationVariables>(CreateSenderRequestFromHttpRequestLogDocument, options);
      }
export type CreateSenderRequestFromHttpRequestLogMutationHookResult = ReturnType<typeof useCreateSenderRequestFromHttpRequestLogMutation>;
export type CreateSenderRequestFromHttpRequestLogMutationResult = Apollo.MutationResult<CreateSenderRequestFromHttpRequestLogMutation>;
export type CreateSenderRequestFromHttpRequestLogMutationOptions = Apollo.BaseMutationOptions<CreateSenderRequestFromHttpRequestLogMutation, CreateSenderRequestFromHttpRequestLogMutationVariables>;
export const SendRequestDocument = gql`
    mutation SendRequest($id: ID!) {
  sendRequest(id: $id) {
    id
  }
}
    `;
export type SendRequestMutationFn = Apollo.MutationFunction<SendRequestMutation, SendRequestMutationVariables>;

/**
 * __useSendRequestMutation__
 *
 * To run a mutation, you first call `useSendRequestMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSendRequestMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [sendRequestMutation, { data, loading, error }] = useSendRequestMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useSendRequestMutation(baseOptions?: Apollo.MutationHookOptions<SendRequestMutation, SendRequestMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SendRequestMutation, SendRequestMutationVariables>(SendRequestDocument, options);
      }
export type SendRequestMutationHookResult = ReturnType<typeof useSendRequestMutation>;
export type SendRequestMutationResult = Apollo.MutationResult<SendRequestMutation>;
export type SendRequestMutationOptions = Apollo.BaseMutationOptions<SendRequestMutation, SendRequestMutationVariables>;
export const GetSenderRequestDocument = gql`
    query GetSenderRequest($id: ID!) {
  senderRequest(id: $id) {
    id
    sourceRequestLogID
    url
    method
    proto
    headers {
      key
      value
    }
    body
    timestamp
    response {
      id
      proto
      statusCode
      statusReason
      body
      headers {
        key
        value
      }
    }
  }
}
    `;

/**
 * __useGetSenderRequestQuery__
 *
 * To run a query within a React component, call `useGetSenderRequestQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSenderRequestQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSenderRequestQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetSenderRequestQuery(baseOptions: Apollo.QueryHookOptions<GetSenderRequestQuery, GetSenderRequestQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSenderRequestQuery, GetSenderRequestQueryVariables>(GetSenderRequestDocument, options);
      }
export function useGetSenderRequestLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSenderRequestQuery, GetSenderRequestQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSenderRequestQuery, GetSenderRequestQueryVariables>(GetSenderRequestDocument, options);
        }
export type GetSenderRequestQueryHookResult = ReturnType<typeof useGetSenderRequestQuery>;
export type GetSenderRequestLazyQueryHookResult = ReturnType<typeof useGetSenderRequestLazyQuery>;
export type GetSenderRequestQueryResult = Apollo.QueryResult<GetSenderRequestQuery, GetSenderRequestQueryVariables>;
export const GetSenderRequestsDocument = gql`
    query GetSenderRequests {
  senderRequests {
    id
    url
    method
    response {
      id
      statusCode
      statusReason
    }
  }
}
    `;

/**
 * __useGetSenderRequestsQuery__
 *
 * To run a query within a React component, call `useGetSenderRequestsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSenderRequestsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSenderRequestsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetSenderRequestsQuery(baseOptions?: Apollo.QueryHookOptions<GetSenderRequestsQuery, GetSenderRequestsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSenderRequestsQuery, GetSenderRequestsQueryVariables>(GetSenderRequestsDocument, options);
      }
export function useGetSenderRequestsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSenderRequestsQuery, GetSenderRequestsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSenderRequestsQuery, GetSenderRequestsQueryVariables>(GetSenderRequestsDocument, options);
        }
export type GetSenderRequestsQueryHookResult = ReturnType<typeof useGetSenderRequestsQuery>;
export type GetSenderRequestsLazyQueryHookResult = ReturnType<typeof useGetSenderRequestsLazyQuery>;
export type GetSenderRequestsQueryResult = Apollo.QueryResult<GetSenderRequestsQuery, GetSenderRequestsQueryVariables>;