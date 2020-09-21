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

  if (loading) return "Loading...";
  if (error) return `Error: ${error.message}`;

  const { method, url, body, response } = data.httpRequestLog;

  return (
    <div>
      <Grid container item spacing={2}>
        <Grid item xs={6}>
          <Box component={Paper} maxHeight="60vh" overflow="scroll">
            <RequestDetail request={{ method, url, body }} />
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
