import {
  makeStyles,
  Theme,
  createStyles,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
} from "@material-ui/core";

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
  };
  return createStyles({
    table: {
      tableLayout: "fixed",
      width: "100%",
    },
    keyCell: {
      ...tableCell,
      width: "40%",
      fontWeight: "bold",
    },
    valueCell: {
      ...tableCell,
      width: "60%",
      border: "none",
      wordBreak: "break-all",
      whiteSpace: "pre-wrap",
    },
  });
});

interface Props {
  headers: Array<{ key: string; value: string }>;
}

function HttpHeadersTable({ headers }: Props): JSX.Element {
  const classes = useStyles();
  return (
    <TableContainer>
      <Table className={classes.table} size="small">
        <TableBody>
          {headers.map(({ key, value }, index) => (
            <TableRow key={index}>
              <TableCell component="th" scope="row" className={classes.keyCell}>
                <code>{key}:</code>
              </TableCell>
              <TableCell className={classes.valueCell}>
                <code>{value}</code>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

export default HttpHeadersTable;
