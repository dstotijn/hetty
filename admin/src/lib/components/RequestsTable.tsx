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
  onRowClick?: (id: string) => void;
  onContextMenu?: (e: React.MouseEvent, id: string) => void;
  rowActions?: (id: string) => JSX.Element;
}

export default function RequestsTable(props: Props): JSX.Element {
  const { requests, activeRowId, onRowClick, onContextMenu, rowActions } = props;

  return (
    <TableContainer sx={{ overflowX: "initial" }}>
      <Table size="small" stickyHeader>
        <TableHead>
          <TableRow>
            <TableCell>Method</TableCell>
            <TableCell>Origin</TableCell>
            <TableCell>Path</TableCell>
            <TableCell>Status</TableCell>
            {rowActions && <TableCell padding="checkbox" />}
          </TableRow>
        </TableHead>
        <TableBody>
          {requests.map(({ id, method, url, response }) => {
            const { origin, pathname, search, hash } = new URL(url);

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
                <PathTableCell>{decodeURIComponent(pathname + search + hash)}</PathTableCell>
                <StatusTableCell>
                  {response && <Status code={response.statusCode} reason={response.statusReason} />}
                </StatusTableCell>
                {rowActions && <TableCell sx={{ py: 0 }}>{rowActions(id)}</TableCell>}
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
