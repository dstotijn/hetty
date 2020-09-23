import React from "react";
import {
  Typography,
  Box,
  createStyles,
  makeStyles,
  Theme,
  Divider,
} from "@material-ui/core";

import HttpHeadersTable from "./HttpHeadersTable";
import Editor from "./Editor";

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
  request: {
    method: string;
    url: string;
    proto: string;
    headers: Array<{ key: string; value: string }>;
    body?: string;
  };
}

function RequestDetail({ request }: Props): JSX.Element {
  const { method, url, proto, headers, body } = request;
  const classes = useStyles();

  const contentType = headers.find((header) => header.key === "Content-Type")
    ?.value;
  const parsedUrl = new URL(url);

  return (
    <div>
      <Box mx={2} my={2}>
        <Typography
          variant="overline"
          color="textSecondary"
          style={{ float: "right" }}
        >
          Request
        </Typography>
        <Typography className={classes.requestTitle} variant="h6">
          {method} {decodeURIComponent(parsedUrl.pathname + parsedUrl.search)}{" "}
          <Typography
            component="span"
            color="textSecondary"
            style={{ fontFamily: "'JetBrains Mono', monospace" }}
          >
            {proto}
          </Typography>
        </Typography>
      </Box>

      <Divider />

      <Box m={2}>
        <HttpHeadersTable headers={headers} />
      </Box>

      {body && <Editor content={body} contentType={contentType} />}
    </div>
  );
}

export default RequestDetail;
