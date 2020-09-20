import { Typography, Box } from "@material-ui/core";
import { green } from "@material-ui/core/colors";
import FiberManualRecordIcon from "@material-ui/icons/FiberManualRecord";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { materialLight } from "react-syntax-highlighter/dist/cjs/styles/prism";

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
        <Typography variant="h5">
          {statusIcon(response.statusCode)} {response.proto} {response.status}
        </Typography>
      </Box>
      <Box>
        <SyntaxHighlighter
          language="markup"
          showLineNumbers={true}
          style={materialLight}
        >
          {response.body}
        </SyntaxHighlighter>
      </Box>
    </div>
  );
}

function statusIcon(status: number): JSX.Element {
  const style = { marginTop: ".2rem", verticalAlign: "top" };
  switch (Math.floor(status / 100)) {
    case 2:
    case 3:
      return <FiberManualRecordIcon style={{ ...style, color: green[400] }} />;
    case 4:
      return <FiberManualRecordIcon style={style} htmlColor={"#f00"} />;
    case 5:
      return <FiberManualRecordIcon style={style} htmlColor={"#f00"} />;
    default:
      return <FiberManualRecordIcon />;
  }
}

export default ResponseDetail;
