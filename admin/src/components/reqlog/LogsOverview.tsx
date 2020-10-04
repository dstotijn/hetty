import { useRouter } from "next/router";
import { gql, useQuery } from "@apollo/client";
import { useState } from "react";
import { Box, Typography, CircularProgress } from "@material-ui/core";
import Alert from "@material-ui/lab/Alert";

import RequestList from "./RequestList";
import LogDetail from "./LogDetail";
import CenteredPaper from "../CenteredPaper";

const HTTP_REQUEST_LOGS = gql`
  query HttpRequestLogs {
    httpRequestLogs {
      id
      method
      url
      timestamp
      response {
        status
        statusCode
      }
    }
  }
`;

function LogsOverview(): JSX.Element {
  const router = useRouter();
  const detailReqLogId =
    router.query.id && parseInt(router.query.id as string, 10);

  const { loading, error, data } = useQuery(HTTP_REQUEST_LOGS, {
    pollInterval: 1000,
  });

  const handleLogClick = (reqId: number) => {
    router.push("/proxy/logs?id=" + reqId, undefined, {
      shallow: false,
    });
  };

  if (loading) {
    return <CircularProgress />;
  }
  if (error) {
    return <Alert severity="error">Error fetching logs: {error.message}</Alert>;
  }

  const { httpRequestLogs: logs } = data;

  return (
    <div>
      <Box mb={2}>
        <RequestList
          logs={logs}
          selectedReqLogId={detailReqLogId}
          onLogClick={handleLogClick}
        />
      </Box>
      <Box>
        {detailReqLogId && <LogDetail requestId={detailReqLogId} />}
        {logs.length !== 0 && !detailReqLogId && (
          <CenteredPaper>
            <Typography>Select a log entry…</Typography>
          </CenteredPaper>
        )}
      </Box>
    </div>
  );
}

export default LogsOverview;
