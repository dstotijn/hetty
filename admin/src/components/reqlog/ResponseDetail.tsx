import { Typography, Box } from "@material-ui/core";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/cjs/styles/prism";

import HttpStatusIcon from "./HttpStatusCode";

interface Props {
  response: {
    proto: string;
    statusCode: number;
    status: string;
    body?: string;
  };
}

function ResponseDetail({ response }: Props): JSX.Element {
  return (
    <div>
      <Box m={3}>
        <Typography
          variant="h6"
          style={{ fontSize: "1rem", whiteSpace: "nowrap" }}
        >
          <HttpStatusIcon status={response.statusCode} /> {response.proto}{" "}
          {response.status}
        </Typography>
      </Box>
      <Box>
        {response.body && (
          <SyntaxHighlighter
            language="markup"
            showLineNumbers={true}
            showInlineLineNumbers={true}
            style={vscDarkPlus}
            lineProps={{
              style: {
                display: "block",
                wordBreak: "break-all",
                whiteSpace: "pre-wrap",
              },
            }}
            wrapLines={true}
          >
            {response.body}
          </SyntaxHighlighter>
        )}
      </Box>
    </div>
  );
}

export default ResponseDetail;
