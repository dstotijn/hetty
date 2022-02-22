import { gql, useMutation } from "@apollo/client";

const CREATE_SENDER_REQUEST = gql`
  mutation CreateSenderRequest($request: SenderRequestInput!) {
    createSenderRequest(request: $request) {
      id
    }
  }
`;

export default function useCreateSenderRequest() {
  return useMutation(CREATE_SENDER_REQUEST);
}
