import { gql, useMutation } from "@apollo/client";
import { HTTP_REQUEST_LOGS } from "./useHttpRequestLogs";

const CLEAR_HTTP_REQUEST_LOG = gql`
  mutation ClearHTTPRequestLog {
    clearHTTPRequestLog {
      success
    }
  }
`;

export function useClearHTTPRequestLog() {
  return useMutation(CLEAR_HTTP_REQUEST_LOG, {
    refetchQueries: [{ query: HTTP_REQUEST_LOGS }],
  });
}
