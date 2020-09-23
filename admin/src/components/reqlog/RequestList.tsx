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
  createStyles,
  makeStyles,
  Theme,
  withTheme,
} from "@material-ui/core";

import HttpStatusIcon from "./HttpStatusCode";
import CenteredPaper from "../CenteredPaper";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    requestTitle: {
      width: "calc(100% - 80px)",
      fontSize: "1rem",
      wordBreak: "break-all",
      whiteSpace: "pre-wrap",
    },
    headersTable: {
      tableLayout: "fixed",
      width: "100%",
    },
    headerKeyCell: {
      verticalAlign: "top",
      width: "30%",
      fontWeight: "bold",
    },
    headerValueCell: {
      width: "70%",
      verticalAlign: "top",
      wordBreak: "break-all",
      whiteSpace: "pre-wrap",
    },
  })
);

interface Props {
  logs: Array<any>;
  selectedReqLogId?: string;
  onLogClick(requestId: string): void;
  theme: Theme;
}

function RequestList({
  logs,
  onLogClick,
  selectedReqLogId,
  theme,
}: Props): JSX.Element {
  return (
    <div>
      <RequestListTable
        onLogClick={onLogClick}
        logs={logs}
        selectedReqLogId={selectedReqLogId}
        theme={theme}
      />
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
  selectedReqLogId?: string;
  onLogClick(requestId: string): void;
  theme: Theme;
}

function RequestListTable({
  logs,
  selectedReqLogId,
  onLogClick,
  theme,
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

            const rowStyle = {
              backgroundColor:
                id === selectedReqLogId
                  ? theme.palette.action.selected
                  : "inherit",
            };

            return (
              <TableRow
                key={id}
                style={rowStyle}
                onClick={() => onLogClick(id)}
              >
                <TableCell style={{ ...cellStyle, width: "100px" }}>
                  <code>{method}</code>
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
                      <code>{response.status}</code>
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

export default withTheme(RequestList);
