import { useState } from "react";
import { Box, Paper, Container, Typography } from "@material-ui/core";

import RequestList from "./RequestList";
import LogDetail from "./LogDetail";

function LogsOverview(): JSX.Element {
  const [detailReqLogId, setDetailReqLogId] = useState<string | null>(null);

  const handleLogClick = (reqId: string) => setDetailReqLogId(reqId);

  return (
    <Box style={{ padding: 8 }}>
      <Box mb={2}>
        <RequestList onLogClick={handleLogClick} />
      </Box>
      <Box>
        {detailReqLogId ? (
          <LogDetail requestId={detailReqLogId} />
        ) : (
          <Paper
            elevation={0}
            style={{
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              height: "60vh",
            }}
          >
            <Typography>Select a log entry…</Typography>
          </Paper>
        )}
      </Box>
    </Box>
  );
}

export default LogsOverview;
