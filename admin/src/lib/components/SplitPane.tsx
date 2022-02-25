import { alpha, styled } from "@mui/material/styles";
import ReactSplitPane, { SplitPaneProps } from "react-split-pane";

const BORDER_WIDTH_FACTOR = 1.75;
const SIZE_FACTOR = 4;
const MARGIN_FACTOR = -1.75;

const SplitPane = styled(ReactSplitPane)<SplitPaneProps>(({ theme }) => ({
  ".Resizer": {
    zIndex: theme.zIndex.mobileStepper,
    boxSizing: "border-box",
    backgroundClip: "padding-box",
    backgroundColor: alpha(theme.palette.grey[400], 0.05),
  },
  ".Resizer:hover": {
    transition: "all 0.5s ease",
    backgroundColor: alpha(theme.palette.primary.main, 1),
  },

  ".Resizer.horizontal": {
    height: theme.spacing(SIZE_FACTOR),
    marginTop: theme.spacing(MARGIN_FACTOR),
    marginBottom: theme.spacing(MARGIN_FACTOR),
    borderTop: `${theme.spacing(BORDER_WIDTH_FACTOR)} solid rgba(255, 255, 255, 0)`,
    borderBottom: `${theme.spacing(BORDER_WIDTH_FACTOR)} solid rgba(255, 255, 255, 0)`,
    borderBottomColor: "rgba(255, 255, 255, 0)",
    cursor: "row-resize",
    width: "100%",
  },

  ".Resizer.vertical": {
    width: theme.spacing(SIZE_FACTOR),
    marginLeft: theme.spacing(MARGIN_FACTOR),
    marginRight: theme.spacing(MARGIN_FACTOR),
    borderLeft: `${theme.spacing(BORDER_WIDTH_FACTOR)} solid rgba(255, 255, 255, 0)`,
    borderRight: `${theme.spacing(BORDER_WIDTH_FACTOR)} solid rgba(255, 255, 255, 0)`,
    cursor: "col-resize",
  },

  ".Resizer.disabled": {
    cursor: "not-allowed",
  },

  ".Resizer.disabled:hover": {
    borderColor: "transparent",
  },

  ".Pane": {
    overflow: "hidden",
  },
}));

export default SplitPane;
