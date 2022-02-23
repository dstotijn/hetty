import { Typography, Box, Divider } from "@mui/material";

import HttpHeadersTable from "./HttpHeadersTable";

import Editor from "lib/components/Editor";
import HttpStatusIcon from "lib/components/HttpStatusIcon";
import { HttpRequestLogQuery } from "lib/graphql/generated";

interface Props {
  response: NonNullable<NonNullable<HttpRequestLogQuery["httpRequestLog"]>["response"]>;
}

function ResponseDetail({ response }: Props): JSX.Element {
  const contentType = response.headers.find((header) => header.key.toLowerCase() === "content-type")?.value;

  return (
    <div>
      <Box p={2}>
        <Typography variant="overline" color="textSecondary" style={{ float: "right" }}>
          Response
        </Typography>
        <Typography variant="h6" style={{ fontSize: "1rem", whiteSpace: "nowrap" }}>
          <HttpStatusIcon status={response.statusCode} />{" "}
          <Typography component="span" color="textSecondary">
            <Typography component="span" color="textSecondary" style={{ fontFamily: "'JetBrains Mono', monospace" }}>
              {response.proto}
            </Typography>
          </Typography>{" "}
          {response.statusCode} {response.statusReason}
        </Typography>
      </Box>

      <Divider />

      <Box p={2}>
        <HttpHeadersTable headers={response.headers} />
      </Box>

      {response.body && <Editor content={response.body} contentType={contentType} />}
    </div>
  );
}

export default ResponseDetail;
