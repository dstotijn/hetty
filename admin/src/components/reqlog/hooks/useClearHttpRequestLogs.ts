import { gql, useMutation } from "@apollo/client";
import { HTTP_REQUEST_LOGS } from "./useHttpRequestLogs";

const CLEAR_REQUEST_LOG = gql`
  mutation ClearRequestLog {
    clearRequestLog {
      success
    }
  }
`;

export function useHttpClearRequestLogs() {
  return useMutation(CLEAR_REQUEST_LOG, {
    refetchQueries: [{ query: HTTP_REQUEST_LOGS }],
  });
}
