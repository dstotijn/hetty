import { Typography, Box } from "@material-ui/core";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { vscDarkPlus } from "react-syntax-highlighter/dist/cjs/styles/prism";

interface Props {
  request: {
    method: string;
    url: string;
    proto: string;
    body?: string;
  };
}

function RequestDetail({ request }: Props): JSX.Element {
  const { method, url, proto, body } = request;

  const parsedUrl = new URL(url);

  return (
    <div>
      <Box m={3}>
        <Typography
          variant="h6"
          style={{ fontSize: "1rem", whiteSpace: "nowrap" }}
        >
          {method} {decodeURIComponent(parsedUrl.pathname + parsedUrl.search)}{" "}
          {proto}
        </Typography>
      </Box>
      <Box>
        {body && (
          <SyntaxHighlighter
            language="markup"
            showLineNumbers={true}
            style={vscDarkPlus}
          >
            {body}
          </SyntaxHighlighter>
        )}
      </Box>
    </div>
  );
}

export default RequestDetail;
