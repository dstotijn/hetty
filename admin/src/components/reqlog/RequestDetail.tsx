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
  return (
    <div>
      <Box m={3}>
        <Typography variant="h5">
          {request.method} {request.url}
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
