import {
  TableContainer,
  Paper,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  CircularProgress,
  Typography,
  Box,
} from "@material-ui/core";

import HttpStatusIcon from "./HttpStatusCode";
import CenteredPaper from "../CenteredPaper";

interface Props {
  logs: Array<any>;
  onLogClick(requestId: string): void;
}

function RequestList({ logs, onLogClick }: Props): JSX.Element {
  return (
    <div>
      <RequestListTable onLogClick={onLogClick} logs={logs} />
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
  logs?: any;
  onLogClick(requestId: string): void;
}

function RequestListTable({
  logs,
  onLogClick,
}: RequestListTableProps): JSX.Element {
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
              <TableRow key={id} onClick={() => onLogClick(id)}>
                <TableCell style={{ ...cellStyle, width: "100px" }}>
                  {method}
                </TableCell>
                <TableCell style={{ ...cellStyle, maxWidth: "100px" }}>
                  {origin}
                </TableCell>
                <TableCell style={{ ...cellStyle, maxWidth: "200px" }}>
                  {decodeURIComponent(pathname + search + hash)}
                </TableCell>
                <TableCell style={{ maxWidth: "100px" }}>
                  {response && (
                    <div>
                      <HttpStatusIcon status={response.statusCode} />{" "}
                      {response.status}
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

export default RequestList;
