import { Alert } from "@mui/lab";
import { Table, TableBody, TableCell, TableContainer, TableRow, Snackbar, SxProps, Theme } from "@mui/material";
import React, { useState } from "react";

const baseCellStyle: SxProps<Theme> = {
  px: 0,
  py: 0.33,
  verticalAlign: "top",
  border: "none",
  whiteSpace: "nowrap",
  overflow: "hidden",
  textOverflow: "ellipsis",
  "&:hover": {
    color: "primary.main",
    whiteSpace: "inherit",
    overflow: "inherit",
    textOverflow: "inherit",
    cursor: "copy",
  },
};

const keyCellStyle = {
  ...baseCellStyle,
  pr: 1,
  width: "40%",
  fontWeight: "bold",
  fontSize: ".75rem",
};

const valueCellStyle = {
  ...baseCellStyle,
  width: "60%",
  border: "none",
  fontSize: ".75rem",
};

interface Props {
  headers: Array<{ key: string; value: string }>;
}

function HttpHeadersTable({ headers }: Props): JSX.Element {
  const [open, setOpen] = useState(false);

  const handleClick = (e: React.MouseEvent) => {
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

    setOpen(true);
  };

  const handleClose = (event: Event | React.SyntheticEvent, reason?: string) => {
    if (reason === "clickaway") {
      return;
    }

    setOpen(false);
  };

  return (
    <div>
      <Snackbar open={open} autoHideDuration={3000} onClose={handleClose}>
        <Alert onClose={handleClose} severity="info">
          Copied to clipboard.
        </Alert>
      </Snackbar>
      <TableContainer>
        <Table
          sx={{
            tableLayout: "fixed",
            width: "100%",
          }}
          size="small"
        >
          <TableBody>
            {headers.map(({ key, value }, index) => (
              <TableRow key={index}>
                <TableCell component="th" scope="row" sx={keyCellStyle} onClick={handleClick}>
                  <code>{key}:</code>
                </TableCell>
                <TableCell sx={valueCellStyle} onClick={handleClick}>
                  <code>{value}</code>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
}

export default HttpHeadersTable;
