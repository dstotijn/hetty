import { teal, orange, red } from "@material-ui/core/colors";
import FiberManualRecordIcon from "@material-ui/icons/FiberManualRecord";

function HttpStatusIcon({ status }: { status: number }): JSX.Element {
  const style = { marginTop: "-.25rem", verticalAlign: "middle" };
  switch (Math.floor(status / 100)) {
    case 2:
    case 3:
      return <FiberManualRecordIcon style={{ ...style, color: teal[400] }} />;
    case 4:
      return <FiberManualRecordIcon style={{ ...style, color: orange[400] }} />;
    case 5:
      return <FiberManualRecordIcon style={{ ...style, color: red[400] }} />;
    default:
      return <FiberManualRecordIcon style={style} />;
  }
}

export default HttpStatusIcon;
