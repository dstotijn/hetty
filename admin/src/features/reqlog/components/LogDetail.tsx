import Alert from "@mui/lab/Alert";
import { Box, CircularProgress, Paper, Typography } from "@mui/material";

import RequestDetail from "./RequestDetail";

import Response from "lib/components/Response";
import SplitPane from "lib/components/SplitPane";
import { useHttpRequestLogQuery } from "lib/graphql/generated";

interface Props {
  id?: string;
}

function LogDetail({ id }: Props): JSX.Element {
  const { loading, error, data } = useHttpRequestLogQuery({
    variables: { id: id as string },
    skip: id === undefined,
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
    return (
      <Paper variant="centered" sx={{ mt: 2 }}>
        <Typography>Select a log entry…</Typography>
      </Paper>
    );
  }

  const reqLog = data.httpRequestLog;

  return (
    <SplitPane split="vertical" size={"50%"}>
      <RequestDetail request={reqLog} />
      {reqLog.response && (
        <Box sx={{ height: "100%", pt: 1, pl: 2, pb: 2 }}>
          <Response response={reqLog.response} />
        </Box>
      )}
    </SplitPane>
  );
}

export default LogDetail;
