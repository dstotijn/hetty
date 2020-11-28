import { useRouter } from "next/router";
import Link from "next/link";
import {
  Box,
  CircularProgress,
  Link as MaterialLink,
  Typography,
} from "@material-ui/core";
import Alert from "@material-ui/lab/Alert";

import RequestList from "./RequestList";
import LogDetail from "./LogDetail";
import CenteredPaper from "../CenteredPaper";
import { useHttpRequestLogs } from "./hooks/useHttpRequestLogs";

function LogsOverview(): JSX.Element {
  const router = useRouter();
  const detailReqLogId =
    router.query.id && parseInt(router.query.id as string, 10);

  const { loading, error, data } = useHttpRequestLogs();

  const handleLogClick = (reqId: number) => {
    router.push("/proxy/logs?id=" + reqId, undefined, {
      shallow: false,
    });
  };

  if (loading) {
    return <CircularProgress />;
  }
  if (error) {
    if (error.graphQLErrors[0]?.extensions?.code === "no_active_project") {
      return (
        <Alert severity="info">
          There is no project active.{" "}
          <Link href="/projects" passHref>
            <MaterialLink color="secondary">Create or open</MaterialLink>
          </Link>{" "}
          one first.
        </Alert>
      );
    }
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
