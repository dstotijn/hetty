import Alert from "@mui/lab/Alert";
import { Box, Grid, Paper, CircularProgress } from "@mui/material";

import RequestDetail from "./RequestDetail";
import ResponseDetail from "./ResponseDetail";

import { useHttpRequestLogQuery } from "lib/graphql/generated";

interface Props {
  requestId: string;
}

function LogDetail({ requestId: id }: Props): JSX.Element {
  const { loading, error, data } = useHttpRequestLogQuery({
    variables: { id },
  });

  if (loading) {
    return <CircularProgress />;
  }
  if (error) {
    return <Alert severity="error">Error fetching logs details: {error.message}</Alert>;
  }

  if (data && !data.httpRequestLog) {
    return (
      <Alert severity="warning">
        Request <strong>{id}</strong> was not found.
      </Alert>
    );
  }

  if (!data?.httpRequestLog) {
    return <div></div>;
  }

  const httpRequestLog = data.httpRequestLog;

  return (
    <div>
      <Grid container item spacing={2}>
        <Grid item xs={6}>
          <Box component={Paper}>
            <RequestDetail request={httpRequestLog} />
          </Box>
        </Grid>
        <Grid item xs={6}>
          {httpRequestLog.response && (
            <Box component={Paper}>
              <ResponseDetail response={httpRequestLog.response} />
            </Box>
          )}
        </Grid>
      </Grid>
    </div>
  );
}

export default LogDetail;
