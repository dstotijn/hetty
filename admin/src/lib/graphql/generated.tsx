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