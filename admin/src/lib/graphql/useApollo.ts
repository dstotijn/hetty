import { ApolloClient, HttpLink, InMemoryCache, NormalizedCacheObject } from "@apollo/client";

let apolloClient: ApolloClient<NormalizedCacheObject>;

function createApolloClient() {
  return new ApolloClient({
    ssrMode: typeof window === "undefined",
    link: new HttpLink({
      uri: "/api/graphql/",
    }),
    cache: new InMemoryCache({
      typePolicies: {
        Query: {
          fields: {
            interceptedRequests: {
              merge(_, incoming) {
                return incoming;
              },
            },
          },
        },
        ProjectSettings: {
          merge: true,
        },
      },
    }),
  });
}

export function useApollo() {
  const _apolloClient = apolloClient ?? createApolloClient();

  // For SSG and SSR always create a new Apollo Client
  if (typeof window === "undefined") return _apolloClient;
  // Create the Apollo Client once in the client
  if (!apolloClient) apolloClient = _apolloClient;

  return _apolloClient;
}
