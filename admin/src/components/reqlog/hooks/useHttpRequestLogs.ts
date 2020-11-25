import { gql, useQuery } from "@apollo/client";

export const HTTP_REQUEST_LOGS = gql`
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

export function useHttpRequestLogs() {
  return useQuery(HTTP_REQUEST_LOGS, {
    pollInterval: 1000,
  });
}
