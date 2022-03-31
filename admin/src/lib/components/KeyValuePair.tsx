import ClearIcon from "@mui/icons-material/Clear";
import {
  Alert,
  IconButton,
  InputBase,
  InputBaseProps,
  Snackbar,
  styled,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TableRowProps,
} from "@mui/material";
import { useState } from "react";

const StyledInputBase = styled(InputBase)<InputBaseProps>(() => ({
  fontSize: "0.875rem",
  "&.MuiInputBase-root input": {
    p: 0,
  },
}));

const StyledTableRow = styled(TableRow)<TableRowProps>(() => ({
  "& .delete-button": {
    visibility: "hidden",
  },
  "&:hover .delete-button": {
    visibility: "inherit",
  },
}));

export interface KeyValuePair {
  key: string;
  value: string;
}

export interface KeyValuePairTableProps {
  items: KeyValuePair[];
  onChange?: (key: string, value: string, index: number) => void;
  onDelete?: (index: number) => void;
}

export function KeyValuePairTable({ items, onChange, onDelete }: KeyValuePairTableProps): JSX.Element {
  const [copyConfOpen, setCopyConfOpen] = useState(false);

  const handleCellClick = (e: React.MouseEvent) => {
    e.preventDefault();

    const windowSel = window.getSelection();

    if (!windowSel || !document) {
      return;
    }

    const r = document.createRange();
    r.selectNode(e.currentTarget);
    windowSel.removeAllRanges();
    windowSel.addRange(r);
    document.execCommand("copy");
    windowSel.removeAllRanges();

    setCopyConfOpen(true);
  };

  const handleCopyConfClose = (_: Event | React.SyntheticEvent, reason?: string) => {
    if (reason === "clickaway") {
      return;
    }

    setCopyConfOpen(false);
  };

  return (
    <div>
      <Snackbar open={copyConfOpen} autoHideDuration={3000} onClose={handleCopyConfClose}>
        <Alert onClose={handleCopyConfClose} severity="info">
          Copied to clipboard.
        </Alert>
      </Snackbar>
      <TableContainer sx={{ overflowX: "initial" }}>
        <Table size="small" stickyHeader>
          <TableHead>
            <TableRow>
              <TableCell>Key</TableCell>
              <TableCell>Value</TableCell>
              {onDelete && <TableCell padding="checkbox"></TableCell>}
            </TableRow>
          </TableHead>
          <TableBody
            sx={{
              "td, th, input": {
                fontFamily: "'JetBrains Mono', monospace",
                fontSize: "0.75rem",
                py: 0.2,
              },
              "td span, th span": {
                display: "block",
                py: 0.7,
              },
            }}
          >
            {items.map(({ key, value }, idx) => (
              <StyledTableRow key={idx} hover>
                <TableCell
                  component="th"
                  scope="row"
                  onClick={(e) => {
                    !onChange && handleCellClick(e);
                  }}
                  sx={{
                    ...(!onChange && {
                      "&:hover": {
                        cursor: "copy",
                      },
                    }),
                  }}
                >
                  {!onChange && <span>{key}</span>}
                  {onChange && (
                    <StyledInputBase
                      size="small"
                      fullWidth
                      placeholder="Key"
                      value={key}
                      onChange={(e) => {
                        onChange && onChange(e.target.value, value, idx);
                      }}
                    />
                  )}
                </TableCell>
                <TableCell
                  onClick={(e) => {
                    !onChange && handleCellClick(e);
                  }}
                  sx={{
                    width: "60%",
                    wordBreak: "break-all",
                    ...(!onChange && {
                      "&:hover": {
                        cursor: "copy",
                      },
                    }),
                  }}
                >
                  {!onChange && value}
                  {onChange && (
                    <StyledInputBase
                      size="small"
                      fullWidth
                      placeholder="Value"
                      value={value}
                      onChange={(e) => {
                        onChange && onChange(key, e.target.value, idx);
                      }}
                    />
                  )}
                </TableCell>
                {onDelete && (
                  <TableCell>
                    <div className="delete-button">
                      <IconButton
                        size="small"
                        onClick={() => {
                          onDelete && onDelete(idx);
                        }}
                        sx={{
                          visibility: onDelete === undefined || items.length === idx + 1 ? "hidden" : "inherit",
                        }}
                      >
                        <ClearIcon fontSize="inherit" />
                      </IconButton>
                    </div>
                  </TableCell>
                )}
              </StyledTableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
}

export default KeyValuePairTable;
