import { Typography, Box } from "@material-ui/core";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { materialLight } from "react-syntax-highlighter/dist/cjs/styles/prism";

interface Props {
  request: {
    method: string;
    url: string;
    body?: string;
  };
}

function RequestDetail({ request }: Props): JSX.Element {
  const { method, url, body } = request;

  const parsedUrl = new URL(url);
  console.log(parsedUrl);

  return (
    <div>
      <Box m={3}>
        <Typography
          variant="h6"
          style={{ fontSize: "1rem", whiteSpace: "nowrap" }}
        >
          {request.method}{" "}
          {decodeURIComponent(parsedUrl.pathname + parsedUrl.search)}
        </Typography>
      </Box>
      <Box>
        {request.body && (
          <SyntaxHighlighter
            language="markup"
            showLineNumbers={true}
            style={materialLight}
          >
            {request.body}
          </SyntaxHighlighter>
        )}
      </Box>
    </div>
  );
}

export default RequestDetail;
