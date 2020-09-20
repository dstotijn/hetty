import { gql, useQuery } from "@apollo/client";
import {
  TableContainer,
  Paper,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  makeStyles,
} from "@material-ui/core";

const HTTP_REQUEST_LOGS = gql`
  query HttpRequestLogs {
    httpRequestLogs {
      id
      method
      url
      timestamp
    }
  }
`;

interface Props {
  onLogClick(requestId: string): void;
}

function RequestList({ onLogClick }: Props): JSX.Element {
  const { loading, error, data } = useQuery(HTTP_REQUEST_LOGS);

  if (loading) return "Loading...";
  if (error) return `Error: ${error.message}`;

  const { httpRequestLogs: logs } = data;

  return (
    <TableContainer component={Paper}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Method</TableCell>
            <TableCell>URL</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {logs.map(({ id, method, url }) => (
            <TableRow key={id} onClick={() => onLogClick(id)}>
              <TableCell>{method}</TableCell>
              <TableCell>{url}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

export default RequestList;
