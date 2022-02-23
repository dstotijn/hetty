import * as colors from "@mui/material/colors";
import { createTheme } from "@mui/material/styles";

const heading = {
  fontFamily: "'JetBrains Mono', monospace",
  fontWeight: 600,
};

let theme = createTheme({
  palette: {
    mode: "dark",
    primary: {
      main: colors.teal["A400"],
    },
    secondary: {
      main: colors.grey[900],
      light: "#333",
      dark: colors.common.black,
    },
  },
  typography: {
    h2: heading,
    h3: heading,
    h4: heading,
    h5: heading,
    h6: heading,
  },
});

theme = createTheme(theme, {
  palette: {
    background: {
      default: theme.palette.secondary.main,
      paper: theme.palette.secondary.light,
    },
    info: {
      main: theme.palette.primary.main,
    },
    success: {
      main: theme.palette.primary.main,
    },
  },
  components: {
    MuiTableCell: {
      styleOverrides: {
        stickyHeader: {
          backgroundColor: theme.palette.secondary.dark,
        },
      },
    },
  },
});

export default theme;
