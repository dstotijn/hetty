import { Typography, Box, Divider } from "@material-ui/core";

import HttpStatusIcon from "./HttpStatusCode";
import Editor from "./Editor";
import HttpHeadersTable from "./HttpHeadersTable";

interface Props {
  response: {
    proto: string;
    statusCode: number;
    status: string;
    headers: Array<{ key: string; value: string }>;
    body?: string;
  };
}

function ResponseDetail({ response }: Props): JSX.Element {
  const contentType = response.headers.find(
    (header) => header.key === "Content-Type"
  )?.value;
  return (
    <div>
      <Box mx={2} my={2}>
        <Typography
          variant="overline"
          color="textSecondary"
          style={{ float: "right" }}
        >
          Response
        </Typography>
        <Typography
          variant="h6"
          style={{ fontSize: "1rem", whiteSpace: "nowrap" }}
        >
          <HttpStatusIcon status={response.statusCode} />{" "}
          <Typography component="span" color="textSecondary">
            {response.proto}
          </Typography>{" "}
          {response.status}
        </Typography>
      </Box>

      <Divider />

      <HttpHeadersTable headers={response.headers} />

      {response.body && (
        <Editor content={response.body} contentType={contentType} />
      )}
    </div>
  );
}

export default ResponseDetail;
