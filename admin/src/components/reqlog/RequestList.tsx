import { gql, useQuery } from "@apollo/client";
import {
  TableContainer,
  Paper,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  Typography,
} from "@material-ui/core";

import HttpStatusIcon from "./HttpStatusCode";

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

interface Props {
  onLogClick(requestId: string): void;
}

function RequestList({ onLogClick }: Props): JSX.Element {
  const { loading, error, data } = useQuery(HTTP_REQUEST_LOGS);

  if (loading) return <div>"Loading..."</div>;
  if (error) return <div>`Error: ${error.message}`</div>;

  const { httpRequestLogs: logs } = data;

  return (
    <TableContainer
      component={Paper}
      style={{ minHeight: 200, height: "24vh" }}
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
