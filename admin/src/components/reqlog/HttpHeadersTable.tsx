import {
  makeStyles,
  Theme,
  createStyles,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
  Snackbar,
} from "@material-ui/core";
import { Alert } from "@material-ui/lab";
import React, { useState } from "react";

const useStyles = makeStyles((theme: Theme) => {
  const paddingX = 0;
  const paddingY = theme.spacing(1) / 3;
  const tableCell = {
    paddingLeft: paddingX,
    paddingRight: paddingX,
    paddingTop: paddingY,
    paddingBottom: paddingY,
    verticalAlign: "top",
    border: "none",
    whiteSpace: "nowrap" as any,
    overflow: "hidden",
    textOverflow: "ellipsis",
    "&:hover": {
      color: theme.palette.secondary.main,
      whiteSpace: "inherit" as any,
      overflow: "inherit",
      textOverflow: "inherit",
      cursor: "copy",
    },
  };
  return createStyles({
    root: {},
    table: {
      tableLayout: "fixed",
      width: "100%",
    },
    keyCell: {
      ...tableCell,
      paddingRight: theme.spacing(1),
      width: "40%",
      fontWeight: "bold",
      fontSize: ".75rem",
    },
    valueCell: {
      ...tableCell,
      width: "60%",
      border: "none",
      fontSize: ".75rem",
    },
  });
});

interface Props {
  headers: Array<{ key: string; value: string }>;
}

function HttpHeadersTable({ headers }: Props): JSX.Element {
  const classes = useStyles();

  const [open, setOpen] = useState(false);

  const handleClick = (e: React.MouseEvent) => {
    e.preventDefault();

    const r = document.createRange();
    r.selectNode(e.currentTarget);
    window.getSelection().removeAllRanges();
    window.getSelection().addRange(r);
    document.execCommand("copy");
    window.getSelection().removeAllRanges();

    setOpen(true);
  };

  const handleClose = (event?: React.SyntheticEvent, reason?: string) => {
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
      <TableContainer className={classes.root}>
        <Table className={classes.table} size="small">
          <TableBody>
            {headers.map(({ key, value }, index) => (
              <TableRow key={index}>
                <TableCell
                  component="th"
                  scope="row"
                  className={classes.keyCell}
                  onClick={handleClick}
                >
                  <code>{key}:</code>
                </TableCell>
                <TableCell className={classes.valueCell} onClick={handleClick}>
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
