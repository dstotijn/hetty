import { gql, useQuery } from "@apollo/client";
import { Box, Grid, Paper } from "@material-ui/core";

import ResponseDetail from "./ResponseDetail";
import RequestDetail from "./RequestDetail";

const HTTP_REQUEST_LOG = gql`
  query HttpRequestLog($id: ID!) {
    httpRequestLog(id: $id) {
      id
      method
      url
      proto
      body
      response {
        proto
        status
        statusCode
        body
      }
    }
  }
`;

interface Props {
  requestId: string;
}

function LogDetail({ requestId: id }: Props): JSX.Element {
  const { loading, error, data } = useQuery(HTTP_REQUEST_LOG, {
    variables: { id },
  });

  if (loading) return <div>"Loading..."</div>;
  if (error) return <div>`Error: ${error.message}`</div>;

  const { method, url, proto, body, response } = data.httpRequestLog;

  return (
    <div>
      <Grid container item spacing={2}>
        <Grid item xs={6}>
          <Box component={Paper} maxHeight="60vh" overflow="scroll">
            <RequestDetail request={{ method, url, proto, body }} />
          </Box>
        </Grid>
        <Grid item xs={6}>
          {response && (
            <Box component={Paper} maxHeight="65vh" overflow="scroll">
              <ResponseDetail response={response} />
            </Box>
          )}
        </Grid>
      </Grid>
    </div>
  );
}

export default LogDetail;
