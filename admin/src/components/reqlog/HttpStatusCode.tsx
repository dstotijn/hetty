import { Theme, withTheme } from "@material-ui/core";
import { orange, red } from "@material-ui/core/colors";
import FiberManualRecordIcon from "@material-ui/icons/FiberManualRecord";

interface Props {
  status: number;
  theme: Theme;
}

function HttpStatusIcon({ status, theme }: Props): JSX.Element {
  const style = { marginTop: "-.25rem", verticalAlign: "middle" };
  switch (Math.floor(status / 100)) {
    case 2:
    case 3:
      return (
        <FiberManualRecordIcon
          style={{ ...style, color: theme.palette.secondary.main }}
        />
      );
    case 4:
      return (
        <FiberManualRecordIcon style={{ ...style, color: orange["A400"] }} />
      );
    case 5:
      return <FiberManualRecordIcon style={{ ...style, color: red["A400"] }} />;
    default:
      return <FiberManualRecordIcon style={style} />;
  }
}

export default withTheme(HttpStatusIcon);
