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
} from "@mui/material";

import HttpStatusIcon from "./HttpStatusCode";
import CenteredPaper from "../CenteredPaper";
import { RequestLog } from "../../lib/requestLogs";

interface Props {
  logs: RequestLog[];
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
  logs: RequestLog[];
  selectedReqLogId?: string;
  onLogClick(requestId: string): void;
}

function RequestListTable({ logs, selectedReqLogId, onLogClick }: RequestListTableProps): JSX.Element {
  const theme = useTheme();

  return (
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
            } as any;

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
              >
                <TableCell style={{ ...cellStyle, width: "100px" }}>
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
  );
}
