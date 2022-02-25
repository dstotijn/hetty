import { Typography, Box } from "@mui/material";
import React from "react";

import RequestTabs from "lib/components/RequestTabs";
import { HttpRequestLogQuery } from "lib/graphql/generated";
import { queryParamsFromURL } from "lib/queryParamsFromURL";

interface Props {
  request: NonNullable<HttpRequestLogQuery["httpRequestLog"]>;
}

function RequestDetail({ request }: Props): JSX.Element {
  const { method, url, headers, body } = request;

  const parsedUrl = new URL(url);

  return (
    <Box sx={{ height: "100%", display: "flex", flexDirection: "column", pr: 2, pb: 2 }}>
      <Box sx={{ p: 2, pb: 0 }}>
        <Typography variant="overline" color="textSecondary" style={{ float: "right" }}>
          Request
        </Typography>
        <Typography
          variant="h6"
          component="h2"
          sx={{
            fontSize: "1rem",
            fontFamily: "'JetBrains Mono', monospace",
            display: "block",
            overflow: "hidden",
            whiteSpace: "nowrap",
            textOverflow: "ellipsis",
            pr: 2,
          }}
        >
          {method} {decodeURIComponent(parsedUrl.pathname + parsedUrl.search)}
        </Typography>
      </Box>

      <Box flex="1 auto" overflow="scroll">
        <RequestTabs headers={headers} queryParams={queryParamsFromURL(url)} body={body} />
      </Box>
    </Box>
  );
}

export default RequestDetail;
