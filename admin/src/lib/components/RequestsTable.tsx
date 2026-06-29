import {
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  styled,
  TableCellProps,
  TableRowProps,
} from "@mui/material";

import HttpStatusIcon from "./HttpStatusIcon";

import { HttpMethod } from "lib/graphql/generated";

const baseCellStyle = {
  whiteSpace: "nowrap",
  overflow: "hidden",
  textOverflow: "ellipsis",
} as const;

const MethodTableCell = styled(TableCell)<TableCellProps>(() => ({
  ...baseCellStyle,
  width: "100px",
}));

const OriginTableCell = styled(TableCell)<TableCellProps>(() => ({
  ...baseCellStyle,
  maxWidth: "100px",
}));

const PathTableCell = styled(TableCell)<TableCellProps>(() => ({
  ...baseCellStyle,
  maxWidth: "200px",
}));

const StatusTableCell = styled(TableCell)<TableCellProps>(() => ({
  ...baseCellStyle,
  width: "100px",
}));

const RequestTableRow = styled(TableRow)<TableRowProps>(() => ({
  "&:hover": {
    cursor: "pointer",
  },
}));

interface HttpRequest {
  id: string;
  url: string;
  method: HttpMethod;
  response?: HttpResponse | null;
}

interface HttpResponse {
  statusCode: number;
  statusReason: string;
  body?: string;
}

interface Props {
  requests: HttpRequest[];
  activeRowId?: string;
  actionsCell?: (id: string) => JSX.Element;
  onRowClick?: (id: string) => void;
  onContextMenu?: (e: React.MouseEvent, id: string) => void;
}

function decodeURLPart(value: string): string {
  try {
    return decodeURIComponent(value);
  } catch {
    return value;
  }
}

function parseURLForDisplay(url: string): { origin: string; path: string } {
  try {
    const { origin, pathname, search, hash } = new URL(url);
    return {
      origin,
      path: decodeURLPart(pathname + search + hash),
    };
  } catch {
    return {
      origin: "",
      path: url,
    };
  }
}

export default function RequestsTable(props: Props): JSX.Element {
  const { requests, activeRowId, actionsCell, onRowClick, onContextMenu } = props;

  return (
    <TableContainer sx={{ overflowX: "initial" }}>
      <Table size="small" stickyHeader>
        <TableHead>
          <TableRow>
            <TableCell>Method</TableCell>
            <TableCell>Origin</TableCell>
            <TableCell>Path</TableCell>
            <TableCell>Status</TableCell>
            {actionsCell && <TableCell padding="checkbox"></TableCell>}
          </TableRow>
        </TableHead>
        <TableBody>
          {requests.map(({ id, method, url, response }) => {
            const { origin, path } = parseURLForDisplay(url);

            return (
              <RequestTableRow
                key={id}
                hover
                selected={id === activeRowId}
                onClick={() => {
                  onRowClick && onRowClick(id);
                }}
                onContextMenu={(e) => {
                  onContextMenu && onContextMenu(e, id);
                }}
              >
                <MethodTableCell>
                  <code>{method}</code>
                </MethodTableCell>
                <OriginTableCell>{origin}</OriginTableCell>
                <PathTableCell>{path}</PathTableCell>
                <StatusTableCell>
                  {response && <Status code={response.statusCode} reason={response.statusReason} />}
                </StatusTableCell>
                {actionsCell && actionsCell(id)}
              </RequestTableRow>
            );
          })}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

function Status({ code, reason }: { code: number; reason: string }): JSX.Element {
  return (
    <div>
      <HttpStatusIcon status={code} />{" "}
      <code>
        {code} {reason}
      </code>
    </div>
  );
}
