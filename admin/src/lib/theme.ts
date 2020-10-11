import { createMuiTheme } from "@material-ui/core/styles";
import grey from "@material-ui/core/colors/grey";
import teal from "@material-ui/core/colors/teal";

const theme = createMuiTheme({
  palette: {
    type: "dark",
    primary: {
      main: grey[900],
    },
    secondary: {
      main: teal["A400"],
    },
    info: {
      main: teal["A400"],
    },
    success: {
      main: teal["A400"],
    },
  },
  typography: {
    h2: {
      fontFamily: "'JetBrains Mono', monospace",
      fontWeight: 600,
    },
    h3: {
      fontFamily: "'JetBrains Mono', monospace",
      fontWeight: 600,
    },
    h4: {
      fontFamily: "'JetBrains Mono', monospace",
      fontWeight: 600,
    },
    h5: {
      fontFamily: "'JetBrains Mono', monospace",
      fontWeight: 600,
    },
    h6: {
      fontFamily: "'JetBrains Mono', monospace",
      fontWeight: 600,
    },
  },
  overrides: {
    MuiTableCell: {
      stickyHeader: {
        backgroundColor: grey[900],
      },
    },
  },
});

export default theme;
