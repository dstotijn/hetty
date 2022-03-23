import React, { createContext, useContext } from "react";

import { GetInterceptedRequestsQuery, useGetInterceptedRequestsQuery } from "./graphql/generated";

const InterceptedRequestsContext = createContext<GetInterceptedRequestsQuery["interceptedRequests"] | null>(null);

interface Props {
  children?: React.ReactNode | undefined;
}

export function InterceptedRequestsProvider({ children }: Props): JSX.Element {
  const { data } = useGetInterceptedRequestsQuery({
    pollInterval: 1000,
  });
  const reqs = data?.interceptedRequests || null;

  return <InterceptedRequestsContext.Provider value={reqs}>{children}</InterceptedRequestsContext.Provider>;
}

export function useInterceptedRequests() {
  return useContext(InterceptedRequestsContext);
}
