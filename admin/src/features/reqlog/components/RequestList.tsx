import {
  TableContainer,
  Paper,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  Typography,
  Box,
  useTheme,
  MenuItem,
  Snackbar,
  Alert,
  Link,
} from "@mui/material";
import React, { useState } from "react";

import CenteredPaper from "lib/components/CenteredPaper";
import HttpStatusIcon from "lib/components/HttpStatusIcon";
import useContextMenu from "lib/components/useContextMenu";
import { HttpRequestLogsQuery, useCreateSenderRequestFromHttpRequestLogMutation } from "lib/graphql/generated";

interface Props {
  logs: NonNullable<HttpRequestLogsQuery["httpRequestLogs"]>;
  selectedReqLogId?: string;
  onLogClick(requestId: string): void;
}

export default function RequestList({ logs, onLogClick, selectedReqLogId }: Props): JSX.Element {
  return (
    <div>
      <RequestListTable onLogClick={onLogClick} logs={logs} selectedReqLogId={selectedReqLogId} />
      {logs.length === 0 && (
        <Box my={1}>
          <CenteredPaper>
            <Typography>No logs found.</Typography>
          </CenteredPaper>
        </Box>
      )}
    </div>
  );
}

interface RequestListTableProps {
  logs: HttpRequestLogsQuery["httpRequestLogs"];
  selectedReqLogId?: string;
  onLogClick(requestId: string): void;
}

function RequestListTable({ logs, selectedReqLogId, onLogClick }: RequestListTableProps): JSX.Element {
  const theme = useTheme();

  const [createSenderReqFromLog] = useCreateSenderRequestFromHttpRequestLogMutation({});

  const [copyToSenderId, setCopyToSenderId] = useState("");
  const [Menu, handleContextMenu, handleContextMenuClose] = useContextMenu();

  const handleCopyToSenderClick = () => {
    createSenderReqFromLog({
      variables: {
        id: copyToSenderId,
      },
      onCompleted({ createSenderRequestFromHttpRequestLog }) {
        const { id } = createSenderRequestFromHttpRequestLog;
        setNewSenderReqId(id);
        setCopiedReqNotifOpen(true);
      },
    });
    handleContextMenuClose();
  };

  const [newSenderReqId, setNewSenderReqId] = useState("");
  const [copiedReqNotifOpen, setCopiedReqNotifOpen] = useState(false);
  const handleCloseCopiedNotif = (_: Event | React.SyntheticEvent, reason?: string) => {
    if (reason === "clickaway") {
      return;
    }
    setCopiedReqNotifOpen(false);
  };

  return (
    <div>
      <Menu>
        <MenuItem onClick={handleCopyToSenderClick}>Copy request to Sender</MenuItem>
      </Menu>
      <Snackbar
        open={copiedReqNotifOpen}
        autoHideDuration={3000}
        onClose={handleCloseCopiedNotif}
        anchorOrigin={{ horizontal: "center", vertical: "bottom" }}
      >
        <Alert onClose={handleCloseCopiedNotif} severity="info">
          Request was copied. <Link href={`/sender?id=${newSenderReqId}`}>Edit in Sender.</Link>
        </Alert>
      </Snackbar>

      <TableContainer
        component={Paper}
        style={{
          minHeight: logs.length ? 200 : 0,
          height: logs.length ? "24vh" : "inherit",
        }}
      >
        <Table stickyHeader size="small">
          <TableHead>
            <TableRow>
              <TableCell>Method</TableCell>
              <TableCell>Origin</TableCell>
              <TableCell>Path</TableCell>
              <TableCell>Status</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {logs.map(({ id, method, url, response }) => {
              const { origin, pathname, search, hash } = new URL(url);

              const cellStyle = {
                whiteSpace: "nowrap",
                overflow: "hidden",
                textOverflow: "ellipsis",
              };

              return (
                <TableRow
                  key={id}
                  sx={{
                    "&:hover": {
                      cursor: "pointer",
                    },
                    ...(id === selectedReqLogId && {
                      bgcolor: theme.palette.action.selected,
                    }),
                  }}
                  hover
                  onClick={() => onLogClick(id)}
                  onContextMenu={(e) => {
                    setCopyToSenderId(id);
                    handleContextMenu(e);
                  }}
                >
                  <TableCell sx={{ ...cellStyle, width: "100px" }}>
                    <code>{method}</code>
                  </TableCell>
                  <TableCell sx={{ ...cellStyle, maxWidth: "100px" }}>{origin}</TableCell>
                  <TableCell sx={{ ...cellStyle, maxWidth: "200px" }}>
                    {decodeURIComponent(pathname + search + hash)}
                  </TableCell>
                  <TableCell style={{ maxWidth: "100px" }}>
                    {response && (
                      <div>
                        <HttpStatusIcon status={response.statusCode} />{" "}
                        <code>
                          {response.statusCode} {response.statusReason}
                        </code>
                      </div>
                    )}
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
}
