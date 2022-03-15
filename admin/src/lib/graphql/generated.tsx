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

export type CancelRequestResult = {
  __typename?: 'CancelRequestResult';
  success: Scalars['Boolean'];
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
  Http10 = 'HTTP10',
  Http11 = 'HTTP11',
  Http20 = 'HTTP20'
}

export type HttpRequest = {
  __typename?: 'HttpRequest';
  body?: Maybe<Scalars['String']>;
  headers: Array<HttpHeader>;
  id: Scalars['ID'];
  method: HttpMethod;
  proto: HttpProtocol;
  url: Scalars['URL'];
};

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

export type InterceptSettings = {
  __typename?: 'InterceptSettings';
  enabled: Scalars['Boolean'];
  requestFilter?: Maybe<Scalars['String']>;
};

export type ModifyRequestInput = {
  body?: InputMaybe<Scalars['String']>;
  headers?: InputMaybe<Array<HttpHeaderInput>>;
  id: Scalars['ID'];
  method: HttpMethod;
  proto: HttpProtocol;
  url: Scalars['URL'];
};

export type ModifyRequestResult = {
  __typename?: 'ModifyRequestResult';
  success: Scalars['Boolean'];
};

export type Mutation = {
  __typename?: 'Mutation';
  cancelRequest: CancelRequestResult;
  clearHTTPRequestLog: ClearHttpRequestLogResult;
  closeProject: CloseProjectResult;
  createOrUpdateSenderRequest: SenderRequest;
  createProject?: Maybe<Project>;
  createSenderRequestFromHttpRequestLog: SenderRequest;
  deleteProject: DeleteProjectResult;
  deleteSenderRequests: DeleteSenderRequestsResult;
  modifyRequest: ModifyRequestResult;
  openProject?: Maybe<Project>;
  sendRequest: SenderRequest;
  setHttpRequestLogFilter?: Maybe<HttpRequestLogFilter>;
  setScope: Array<ScopeRule>;
  setSenderRequestFilter?: Maybe<SenderRequestFilter>;
  updateInterceptSettings: InterceptSettings;
};


export type MutationCancelRequestArgs = {
  id: Scalars['ID'];
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


export type MutationModifyRequestArgs = {
  request: ModifyRequestInput;
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


export type MutationUpdateInterceptSettingsArgs = {
  input: UpdateInterceptSettingsInput;
};

export type Project = {
  __typename?: 'Project';
  id: Scalars['ID'];
  isActive: Scalars['Boolean'];
  name: Scalars['String'];
  settings: ProjectSettings;
};

export type ProjectSettings = {
  __typename?: 'ProjectSettings';
  intercept: InterceptSettings;
};

export type Query = {
  __typename?: 'Query';
  activeProject?: Maybe<Project>;
  httpRequestLog?: Maybe<HttpRequestLog>;
  httpRequestLogFilter?: Maybe<HttpRequestLogFilter>;
  httpRequestLogs: Array<HttpRequestLog>;
  interceptedRequest?: Maybe<HttpRequest>;
  interceptedRequests: Array<HttpRequest>;
  projects: Array<Project>;
  scope: Array<ScopeRule>;
  senderRequest?: Maybe<SenderRequest>;
  senderRequests: Array<SenderRequest>;
};


export type QueryHttpRequestLogArgs = {
  id: Scalars['ID'];
};


export type QueryInterceptedRequestArgs = {
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

export type UpdateInterceptSettingsInput = {
  enabled: Scalars['Boolean'];
  requestFilter?: InputMaybe<Scalars['String']>;
};

export type CancelRequestMutationVariables = Exact<{
  id: Scalars['ID'];
}>;


export type CancelRequestMutation = { __typename?: 'Mutation', cancelRequest: { __typename?: 'CancelRequestResult', success: boolean } };

export type GetInterceptedRequestQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetInterceptedRequestQuery = { __typename?: 'Query', interceptedRequest?: { __typename?: 'HttpRequest', id: string, url: any, method: HttpMethod, proto: HttpProtocol, body?: string | null, headers: Array<{ __typename?: 'HttpHeader', key: string, value: string }> } | null };

export type ModifyRequestMutationVariables = Exact<{
  request: ModifyRequestInput;
}>;


export type ModifyRequestMutation = { __typename?: 'Mutation', modifyRequest: { __typename?: 'ModifyRequestResult', success: boolean } };

export type ActiveProjectQueryVariables = Exact<{ [key: string]: never; }>;


export type ActiveProjectQuery = { __typename?: 'Query', activeProject?: { __typename?: 'Project', id: string, name: string, isActive: boolean, settings: { __typename?: 'ProjectSettings', intercept: { __typename?: 'InterceptSettings', enabled: boolean, requestFilter?: string | null } } } | null };

export type CloseProjectMutationVariables = Exact<{ [key: string]: never; }>;


export type CloseProjectMutation = { __typename?: 'Mutation', closeProject: { __typename?: 'CloseProjectResult', success: boolean } };

export type CreateProjectMutationVariables = Exact<{
  name: Scalars['String'];
}>;


export type CreateProjectMutation = { __typename?: 'Mutation', createProject?: { __typename?: 'Project', id: string, name: string } | null };

export type DeleteProjectMutationVariables = Exact<{
  id: Scalars['ID'];
}>;


export type DeleteProjectMutation = { __typename?: 'Mutation', deleteProject: { __typename?: 'DeleteProjectResult', success: boolean } };

export type OpenProjectMutationVariables = Exact<{
  id: Scalars['ID'];
}>;


export type OpenProjectMutation = { __typename?: 'Mutation', openProject?: { __typename?: 'Project', id: string, name: string, isActive: boolean } | null };

export type ProjectsQueryVariables = Exact<{ [key: string]: never; }>;


export type ProjectsQuery = { __typename?: 'Query', projects: Array<{ __typename?: 'Project', id: string, name: string, isActive: boolean }> };

export type ClearHttpRequestLogMutationVariables = Exact<{ [key: string]: never; }>;


export type ClearHttpRequestLogMutation = { __typename?: 'Mutation', clearHTTPRequestLog: { __typename?: 'ClearHTTPRequestLogResult', success: boolean } };

export type HttpRequestLogQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type HttpRequestLogQuery = { __typename?: 'Query', httpRequestLog?: { __typename?: 'HttpRequestLog', id: string, method: HttpMethod, url: string, proto: string, body?: string | null, headers: Array<{ __typename?: 'HttpHeader', key: string, value: string }>, response?: { __typename?: 'HttpResponseLog', id: string, proto: HttpProtocol, statusCode: number, statusReason: string, body?: string | null, headers: Array<{ __typename?: 'HttpHeader', key: string, value: string }> } | null } | null };

export type HttpRequestLogFilterQueryVariables = Exact<{ [key: string]: never; }>;


export type HttpRequestLogFilterQuery = { __typename?: 'Query', httpRequestLogFilter?: { __typename?: 'HttpRequestLogFilter', onlyInScope: boolean, searchExpression?: string | null } | null };

export type HttpRequestLogsQueryVariables = Exact<{ [key: string]: never; }>;


export type HttpRequestLogsQuery = { __typename?: 'Query', httpRequestLogs: Array<{ __typename?: 'HttpRequestLog', id: string, method: HttpMethod, url: string, timestamp: any, response?: { __typename?: 'HttpResponseLog', statusCode: number, statusReason: string } | null }> };

export type SetHttpRequestLogFilterMutationVariables = Exact<{
  filter?: InputMaybe<HttpRequestLogFilterInput>;
}>;


export type SetHttpRequestLogFilterMutation = { __typename?: 'Mutation', setHttpRequestLogFilter?: { __typename?: 'HttpRequestLogFilter', onlyInScope: boolean, searchExpression?: string | null } | null };

export type ScopeQueryVariables = Exact<{ [key: string]: never; }>;


export type ScopeQuery = { __typename?: 'Query', scope: Array<{ __typename?: 'ScopeRule', url?: any | null }> };

export type SetScopeMutationVariables = Exact<{
  scope: Array<ScopeRuleInput> | ScopeRuleInput;
}>;


export type SetScopeMutation = { __typename?: 'Mutation', setScope: Array<{ __typename?: 'ScopeRule', url?: any | null }> };

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

export type UpdateInterceptSettingsMutationVariables = Exact<{
  input: UpdateInterceptSettingsInput;
}>;


export type UpdateInterceptSettingsMutation = { __typename?: 'Mutation', updateInterceptSettings: { __typename?: 'InterceptSettings', enabled: boolean, requestFilter?: string | null } };

export type GetInterceptedRequestsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetInterceptedRequestsQuery = { __typename?: 'Query', interceptedRequests: Array<{ __typename?: 'HttpRequest', id: string, url: any, method: HttpMethod }> };


export const CancelRequestDocument = gql`
    mutation CancelRequest($id: ID!) {
  cancelRequest(id: $id) {
    success
  }
}
    `;
export type CancelRequestMutationFn = Apollo.MutationFunction<CancelRequestMutation, CancelRequestMutationVariables>;

/**
 * __useCancelRequestMutation__
 *
 * To run a mutation, you first call `useCancelRequestMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCancelRequestMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [cancelRequestMutation, { data, loading, error }] = useCancelRequestMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useCancelRequestMutation(baseOptions?: Apollo.MutationHookOptions<CancelRequestMutation, CancelRequestMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CancelRequestMutation, CancelRequestMutationVariables>(CancelRequestDocument, options);
      }
export type CancelRequestMutationHookResult = ReturnType<typeof useCancelRequestMutation>;
export type CancelRequestMutationResult = Apollo.MutationResult<CancelRequestMutation>;
export type CancelRequestMutationOptions = Apollo.BaseMutationOptions<CancelRequestMutation, CancelRequestMutationVariables>;
export const GetInterceptedRequestDocument = gql`
    query GetInterceptedRequest($id: ID!) {
  interceptedRequest(id: $id) {
    id
    url
    method
    proto
    headers {
      key
      value
    }
    body
  }
}
    `;

/**
 * __useGetInterceptedRequestQuery__
 *
 * To run a query within a React component, call `useGetInterceptedRequestQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetInterceptedRequestQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetInterceptedRequestQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetInterceptedRequestQuery(baseOptions: Apollo.QueryHookOptions<GetInterceptedRequestQuery, GetInterceptedRequestQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetInterceptedRequestQuery, GetInterceptedRequestQueryVariables>(GetInterceptedRequestDocument, options);
      }
export function useGetInterceptedRequestLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetInterceptedRequestQuery, GetInterceptedRequestQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetInterceptedRequestQuery, GetInterceptedRequestQueryVariables>(GetInterceptedRequestDocument, options);
        }
export type GetInterceptedRequestQueryHookResult = ReturnType<typeof useGetInterceptedRequestQuery>;
export type GetInterceptedRequestLazyQueryHookResult = ReturnType<typeof useGetInterceptedRequestLazyQuery>;
export type GetInterceptedRequestQueryResult = Apollo.QueryResult<GetInterceptedRequestQuery, GetInterceptedRequestQueryVariables>;
export const ModifyRequestDocument = gql`
    mutation ModifyRequest($request: ModifyRequestInput!) {
  modifyRequest(request: $request) {
    success
  }
}
    `;
export type ModifyRequestMutationFn = Apollo.MutationFunction<ModifyRequestMutation, ModifyRequestMutationVariables>;

/**
 * __useModifyRequestMutation__
 *
 * To run a mutation, you first call `useModifyRequestMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useModifyRequestMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [modifyRequestMutation, { data, loading, error }] = useModifyRequestMutation({
 *   variables: {
 *      request: // value for 'request'
 *   },
 * });
 */
export function useModifyRequestMutation(baseOptions?: Apollo.MutationHookOptions<ModifyRequestMutation, ModifyRequestMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ModifyRequestMutation, ModifyRequestMutationVariables>(ModifyRequestDocument, options);
      }
export type ModifyRequestMutationHookResult = ReturnType<typeof useModifyRequestMutation>;
export type ModifyRequestMutationResult = Apollo.MutationResult<ModifyRequestMutation>;
export type ModifyRequestMutationOptions = Apollo.BaseMutationOptions<ModifyRequestMutation, ModifyRequestMutationVariables>;
export const ActiveProjectDocument = gql`
    query ActiveProject {
  activeProject {
    id
    name
    isActive
    settings {
      intercept {
        enabled
        requestFilter
      }
    }
  }
}
    `;

/**
 * __useActiveProjectQuery__
 *
 * To run a query within a React component, call `useActiveProjectQuery` and pass it any options that fit your needs.
 * When your component renders, `useActiveProjectQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useActiveProjectQuery({
 *   variables: {
 *   },
 * });
 */
export function useActiveProjectQuery(baseOptions?: Apollo.QueryHookOptions<ActiveProjectQuery, ActiveProjectQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ActiveProjectQuery, ActiveProjectQueryVariables>(ActiveProjectDocument, options);
      }
export function useActiveProjectLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ActiveProjectQuery, ActiveProjectQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ActiveProjectQuery, ActiveProjectQueryVariables>(ActiveProjectDocument, options);
        }
export type ActiveProjectQueryHookResult = ReturnType<typeof useActiveProjectQuery>;
export type ActiveProjectLazyQueryHookResult = ReturnType<typeof useActiveProjectLazyQuery>;
export type ActiveProjectQueryResult = Apollo.QueryResult<ActiveProjectQuery, ActiveProjectQueryVariables>;
export const CloseProjectDocument = gql`
    mutation CloseProject {
  closeProject {
    success
  }
}
    `;
export type CloseProjectMutationFn = Apollo.MutationFunction<CloseProjectMutation, CloseProjectMutationVariables>;

/**
 * __useCloseProjectMutation__
 *
 * To run a mutation, you first call `useCloseProjectMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCloseProjectMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [closeProjectMutation, { data, loading, error }] = useCloseProjectMutation({
 *   variables: {
 *   },
 * });
 */
export function useCloseProjectMutation(baseOptions?: Apollo.MutationHookOptions<CloseProjectMutation, CloseProjectMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CloseProjectMutation, CloseProjectMutationVariables>(CloseProjectDocument, options);
      }
export type CloseProjectMutationHookResult = ReturnType<typeof useCloseProjectMutation>;
export type CloseProjectMutationResult = Apollo.MutationResult<CloseProjectMutation>;
export type CloseProjectMutationOptions = Apollo.BaseMutationOptions<CloseProjectMutation, CloseProjectMutationVariables>;
export const CreateProjectDocument = gql`
    mutation CreateProject($name: String!) {
  createProject(name: $name) {
    id
    name
  }
}
    `;
export type CreateProjectMutationFn = Apollo.MutationFunction<CreateProjectMutation, CreateProjectMutationVariables>;

/**
 * __useCreateProjectMutation__
 *
 * To run a mutation, you first call `useCreateProjectMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateProjectMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createProjectMutation, { data, loading, error }] = useCreateProjectMutation({
 *   variables: {
 *      name: // value for 'name'
 *   },
 * });
 */
export function useCreateProjectMutation(baseOptions?: Apollo.MutationHookOptions<CreateProjectMutation, CreateProjectMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateProjectMutation, CreateProjectMutationVariables>(CreateProjectDocument, options);
      }
export type CreateProjectMutationHookResult = ReturnType<typeof useCreateProjectMutation>;
export type CreateProjectMutationResult = Apollo.MutationResult<CreateProjectMutation>;
export type CreateProjectMutationOptions = Apollo.BaseMutationOptions<CreateProjectMutation, CreateProjectMutationVariables>;
export const DeleteProjectDocument = gql`
    mutation DeleteProject($id: ID!) {
  deleteProject(id: $id) {
    success
  }
}
    `;
export type DeleteProjectMutationFn = Apollo.MutationFunction<DeleteProjectMutation, DeleteProjectMutationVariables>;

/**
 * __useDeleteProjectMutation__
 *
 * To run a mutation, you first call `useDeleteProjectMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteProjectMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteProjectMutation, { data, loading, error }] = useDeleteProjectMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useDeleteProjectMutation(baseOptions?: Apollo.MutationHookOptions<DeleteProjectMutation, DeleteProjectMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteProjectMutation, DeleteProjectMutationVariables>(DeleteProjectDocument, options);
      }
export type DeleteProjectMutationHookResult = ReturnType<typeof useDeleteProjectMutation>;
export type DeleteProjectMutationResult = Apollo.MutationResult<DeleteProjectMutation>;
export type DeleteProjectMutationOptions = Apollo.BaseMutationOptions<DeleteProjectMutation, DeleteProjectMutationVariables>;
export const OpenProjectDocument = gql`
    mutation OpenProject($id: ID!) {
  openProject(id: $id) {
    id
    name
    isActive
  }
}
    `;
export type OpenProjectMutationFn = Apollo.MutationFunction<OpenProjectMutation, OpenProjectMutationVariables>;

/**
 * __useOpenProjectMutation__
 *
 * To run a mutation, you first call `useOpenProjectMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useOpenProjectMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [openProjectMutation, { data, loading, error }] = useOpenProjectMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useOpenProjectMutation(baseOptions?: Apollo.MutationHookOptions<OpenProjectMutation, OpenProjectMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<OpenProjectMutation, OpenProjectMutationVariables>(OpenProjectDocument, options);
      }
export type OpenProjectMutationHookResult = ReturnType<typeof useOpenProjectMutation>;
export type OpenProjectMutationResult = Apollo.MutationResult<OpenProjectMutation>;
export type OpenProjectMutationOptions = Apollo.BaseMutationOptions<OpenProjectMutation, OpenProjectMutationVariables>;
export const ProjectsDocument = gql`
    query Projects {
  projects {
    id
    name
    isActive
  }
}
    `;

/**
 * __useProjectsQuery__
 *
 * To run a query within a React component, call `useProjectsQuery` and pass it any options that fit your needs.
 * When your component renders, `useProjectsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useProjectsQuery({
 *   variables: {
 *   },
 * });
 */
export function useProjectsQuery(baseOptions?: Apollo.QueryHookOptions<ProjectsQuery, ProjectsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ProjectsQuery, ProjectsQueryVariables>(ProjectsDocument, options);
      }
export function useProjectsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ProjectsQuery, ProjectsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ProjectsQuery, ProjectsQueryVariables>(ProjectsDocument, options);
        }
export type ProjectsQueryHookResult = ReturnType<typeof useProjectsQuery>;
export type ProjectsLazyQueryHookResult = ReturnType<typeof useProjectsLazyQuery>;
export type ProjectsQueryResult = Apollo.QueryResult<ProjectsQuery, ProjectsQueryVariables>;
export const ClearHttpRequestLogDocument = gql`
    mutation ClearHTTPRequestLog {
  clearHTTPRequestLog {
    success
  }
}
    `;
export type ClearHttpRequestLogMutationFn = Apollo.MutationFunction<ClearHttpRequestLogMutation, ClearHttpRequestLogMutationVariables>;

/**
 * __useClearHttpRequestLogMutation__
 *
 * To run a mutation, you first call `useClearHttpRequestLogMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useClearHttpRequestLogMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [clearHttpRequestLogMutation, { data, loading, error }] = useClearHttpRequestLogMutation({
 *   variables: {
 *   },
 * });
 */
export function useClearHttpRequestLogMutation(baseOptions?: Apollo.MutationHookOptions<ClearHttpRequestLogMutation, ClearHttpRequestLogMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ClearHttpRequestLogMutation, ClearHttpRequestLogMutationVariables>(ClearHttpRequestLogDocument, options);
      }
export type ClearHttpRequestLogMutationHookResult = ReturnType<typeof useClearHttpRequestLogMutation>;
export type ClearHttpRequestLogMutationResult = Apollo.MutationResult<ClearHttpRequestLogMutation>;
export type ClearHttpRequestLogMutationOptions = Apollo.BaseMutationOptions<ClearHttpRequestLogMutation, ClearHttpRequestLogMutationVariables>;
export const HttpRequestLogDocument = gql`
    query HttpRequestLog($id: ID!) {
  httpRequestLog(id: $id) {
    id
    method
    url
    proto
    headers {
      key
      value
    }
    body
    response {
      id
      proto
      headers {
        key
        value
      }
      statusCode
      statusReason
      body
    }
  }
}
    `;

/**
 * __useHttpRequestLogQuery__
 *
 * To run a query within a React component, call `useHttpRequestLogQuery` and pass it any options that fit your needs.
 * When your component renders, `useHttpRequestLogQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useHttpRequestLogQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useHttpRequestLogQuery(baseOptions: Apollo.QueryHookOptions<HttpRequestLogQuery, HttpRequestLogQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<HttpRequestLogQuery, HttpRequestLogQueryVariables>(HttpRequestLogDocument, options);
      }
export function useHttpRequestLogLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<HttpRequestLogQuery, HttpRequestLogQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<HttpRequestLogQuery, HttpRequestLogQueryVariables>(HttpRequestLogDocument, options);
        }
export type HttpRequestLogQueryHookResult = ReturnType<typeof useHttpRequestLogQuery>;
export type HttpRequestLogLazyQueryHookResult = ReturnType<typeof useHttpRequestLogLazyQuery>;
export type HttpRequestLogQueryResult = Apollo.QueryResult<HttpRequestLogQuery, HttpRequestLogQueryVariables>;
export const HttpRequestLogFilterDocument = gql`
    query HttpRequestLogFilter {
  httpRequestLogFilter {
    onlyInScope
    searchExpression
  }
}
    `;

/**
 * __useHttpRequestLogFilterQuery__
 *
 * To run a query within a React component, call `useHttpRequestLogFilterQuery` and pass it any options that fit your needs.
 * When your component renders, `useHttpRequestLogFilterQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useHttpRequestLogFilterQuery({
 *   variables: {
 *   },
 * });
 */
export function useHttpRequestLogFilterQuery(baseOptions?: Apollo.QueryHookOptions<HttpRequestLogFilterQuery, HttpRequestLogFilterQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<HttpRequestLogFilterQuery, HttpRequestLogFilterQueryVariables>(HttpRequestLogFilterDocument, options);
      }
export function useHttpRequestLogFilterLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<HttpRequestLogFilterQuery, HttpRequestLogFilterQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<HttpRequestLogFilterQuery, HttpRequestLogFilterQueryVariables>(HttpRequestLogFilterDocument, options);
        }
export type HttpRequestLogFilterQueryHookResult = ReturnType<typeof useHttpRequestLogFilterQuery>;
export type HttpRequestLogFilterLazyQueryHookResult = ReturnType<typeof useHttpRequestLogFilterLazyQuery>;
export type HttpRequestLogFilterQueryResult = Apollo.QueryResult<HttpRequestLogFilterQuery, HttpRequestLogFilterQueryVariables>;
export const HttpRequestLogsDocument = gql`
    query HttpRequestLogs {
  httpRequestLogs {
    id
    method
    url
    timestamp
    response {
      statusCode
      statusReason
    }
  }
}
    `;

/**
 * __useHttpRequestLogsQuery__
 *
 * To run a query within a React component, call `useHttpRequestLogsQuery` and pass it any options that fit your needs.
 * When your component renders, `useHttpRequestLogsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useHttpRequestLogsQuery({
 *   variables: {
 *   },
 * });
 */
export function useHttpRequestLogsQuery(baseOptions?: Apollo.QueryHookOptions<HttpRequestLogsQuery, HttpRequestLogsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<HttpRequestLogsQuery, HttpRequestLogsQueryVariables>(HttpRequestLogsDocument, options);
      }
export function useHttpRequestLogsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<HttpRequestLogsQuery, HttpRequestLogsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<HttpRequestLogsQuery, HttpRequestLogsQueryVariables>(HttpRequestLogsDocument, options);
        }
export type HttpRequestLogsQueryHookResult = ReturnType<typeof useHttpRequestLogsQuery>;
export type HttpRequestLogsLazyQueryHookResult = ReturnType<typeof useHttpRequestLogsLazyQuery>;
export type HttpRequestLogsQueryResult = Apollo.QueryResult<HttpRequestLogsQuery, HttpRequestLogsQueryVariables>;
export const SetHttpRequestLogFilterDocument = gql`
    mutation SetHttpRequestLogFilter($filter: HttpRequestLogFilterInput) {
  setHttpRequestLogFilter(filter: $filter) {
    onlyInScope
    searchExpression
  }
}
    `;
export type SetHttpRequestLogFilterMutationFn = Apollo.MutationFunction<SetHttpRequestLogFilterMutation, SetHttpRequestLogFilterMutationVariables>;

/**
 * __useSetHttpRequestLogFilterMutation__
 *
 * To run a mutation, you first call `useSetHttpRequestLogFilterMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSetHttpRequestLogFilterMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [setHttpRequestLogFilterMutation, { data, loading, error }] = useSetHttpRequestLogFilterMutation({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useSetHttpRequestLogFilterMutation(baseOptions?: Apollo.MutationHookOptions<SetHttpRequestLogFilterMutation, SetHttpRequestLogFilterMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SetHttpRequestLogFilterMutation, SetHttpRequestLogFilterMutationVariables>(SetHttpRequestLogFilterDocument, options);
      }
export type SetHttpRequestLogFilterMutationHookResult = ReturnType<typeof useSetHttpRequestLogFilterMutation>;
export type SetHttpRequestLogFilterMutationResult = Apollo.MutationResult<SetHttpRequestLogFilterMutation>;
export type SetHttpRequestLogFilterMutationOptions = Apollo.BaseMutationOptions<SetHttpRequestLogFilterMutation, SetHttpRequestLogFilterMutationVariables>;
export const ScopeDocument = gql`
    query Scope {
  scope {
    url
  }
}
    `;

/**
 * __useScopeQuery__
 *
 * To run a query within a React component, call `useScopeQuery` and pass it any options that fit your needs.
 * When your component renders, `useScopeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useScopeQuery({
 *   variables: {
 *   },
 * });
 */
export function useScopeQuery(baseOptions?: Apollo.QueryHookOptions<ScopeQuery, ScopeQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ScopeQuery, ScopeQueryVariables>(ScopeDocument, options);
      }
export function useScopeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ScopeQuery, ScopeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ScopeQuery, ScopeQueryVariables>(ScopeDocument, options);
        }
export type ScopeQueryHookResult = ReturnType<typeof useScopeQuery>;
export type ScopeLazyQueryHookResult = ReturnType<typeof useScopeLazyQuery>;
export type ScopeQueryResult = Apollo.QueryResult<ScopeQuery, ScopeQueryVariables>;
export const SetScopeDocument = gql`
    mutation SetScope($scope: [ScopeRuleInput!]!) {
  setScope(scope: $scope) {
    url
  }
}
    `;
export type SetScopeMutationFn = Apollo.MutationFunction<SetScopeMutation, SetScopeMutationVariables>;

/**
 * __useSetScopeMutation__
 *
 * To run a mutation, you first call `useSetScopeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSetScopeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [setScopeMutation, { data, loading, error }] = useSetScopeMutation({
 *   variables: {
 *      scope: // value for 'scope'
 *   },
 * });
 */
export function useSetScopeMutation(baseOptions?: Apollo.MutationHookOptions<SetScopeMutation, SetScopeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SetScopeMutation, SetScopeMutationVariables>(SetScopeDocument, options);
      }
export type SetScopeMutationHookResult = ReturnType<typeof useSetScopeMutation>;
export type SetScopeMutationResult = Apollo.MutationResult<SetScopeMutation>;
export type SetScopeMutationOptions = Apollo.BaseMutationOptions<SetScopeMutation, SetScopeMutationVariables>;
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
export const UpdateInterceptSettingsDocument = gql`
    mutation UpdateInterceptSettings($input: UpdateInterceptSettingsInput!) {
  updateInterceptSettings(input: $input) {
    enabled
    requestFilter
  }
}
    `;
export type UpdateInterceptSettingsMutationFn = Apollo.MutationFunction<UpdateInterceptSettingsMutation, UpdateInterceptSettingsMutationVariables>;

/**
 * __useUpdateInterceptSettingsMutation__
 *
 * To run a mutation, you first call `useUpdateInterceptSettingsMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateInterceptSettingsMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateInterceptSettingsMutation, { data, loading, error }] = useUpdateInterceptSettingsMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateInterceptSettingsMutation(baseOptions?: Apollo.MutationHookOptions<UpdateInterceptSettingsMutation, UpdateInterceptSettingsMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateInterceptSettingsMutation, UpdateInterceptSettingsMutationVariables>(UpdateInterceptSettingsDocument, options);
      }
export type UpdateInterceptSettingsMutationHookResult = ReturnType<typeof useUpdateInterceptSettingsMutation>;
export type UpdateInterceptSettingsMutationResult = Apollo.MutationResult<UpdateInterceptSettingsMutation>;
export type UpdateInterceptSettingsMutationOptions = Apollo.BaseMutationOptions<UpdateInterceptSettingsMutation, UpdateInterceptSettingsMutationVariables>;
export const GetInterceptedRequestsDocument = gql`
    query GetInterceptedRequests {
  interceptedRequests {
    id
    url
    method
  }
}
    `;

/**
 * __useGetInterceptedRequestsQuery__
 *
 * To run a query within a React component, call `useGetInterceptedRequestsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetInterceptedRequestsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetInterceptedRequestsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetInterceptedRequestsQuery(baseOptions?: Apollo.QueryHookOptions<GetInterceptedRequestsQuery, GetInterceptedRequestsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetInterceptedRequestsQuery, GetInterceptedRequestsQueryVariables>(GetInterceptedRequestsDocument, options);
      }
export function useGetInterceptedRequestsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetInterceptedRequestsQuery, GetInterceptedRequestsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetInterceptedRequestsQuery, GetInterceptedRequestsQueryVariables>(GetInterceptedRequestsDocument, options);
        }
export type GetInterceptedRequestsQueryHookResult = ReturnType<typeof useGetInterceptedRequestsQuery>;
export type GetInterceptedRequestsLazyQueryHookResult = ReturnType<typeof useGetInterceptedRequestsLazyQuery>;
export type GetInterceptedRequestsQueryResult = Apollo.QueryResult<GetInterceptedRequestsQuery, GetInterceptedRequestsQueryVariables>;