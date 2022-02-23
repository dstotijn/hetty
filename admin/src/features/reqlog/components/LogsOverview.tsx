import Alert from "@mui/lab/Alert";
import { Box, CircularProgress, Link as MaterialLink, Typography } from "@mui/material";
import Link from "next/link";
import { useRouter } from "next/router";

import LogDetail from "./LogDetail";
import RequestList from "./RequestList";

import CenteredPaper from "lib/components/CenteredPaper";
import { useHttpRequestLogsQuery } from "lib/graphql/generated";

export default function LogsOverview(): JSX.Element {
  const router = useRouter();
  const detailReqLogId = router.query.id as string | undefined;
  const { loading, error, data } = useHttpRequestLogsQuery({
    pollInterval: 1000,
  });

  const handleLogClick = (reqId: string) => {
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
            <MaterialLink color="primary">Create or open</MaterialLink>
          </Link>{" "}
          one first.
        </Alert>
      );
    }
    return <Alert severity="error">Error fetching logs: {error.message}</Alert>;
  }

  const logs = data?.httpRequestLogs || [];

  return (
    <div>
      <Box mb={2}>
        <RequestList logs={logs} selectedReqLogId={detailReqLogId} onLogClick={handleLogClick} />
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
