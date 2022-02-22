import { SvgIconTypeMap } from "@mui/material";
import FiberManualRecordIcon from "@mui/icons-material/FiberManualRecord";

interface Props {
  status: number;
}

export default function HttpStatusIcon({ status }: Props): JSX.Element {
  let color: SvgIconTypeMap["props"]["color"] = "inherit";

  switch (Math.floor(status / 100)) {
    case 2:
    case 3:
      color = "primary";
      break;
    case 4:
      color = "warning";
      break;
    case 5:
      color = "error";
      break;
  }

  return <FiberManualRecordIcon sx={{ marginTop: "-.25rem", verticalAlign: "middle" }} color={color} />;
}
