import * as colors from "@mui/material/colors";
import { createTheme } from "@mui/material/styles";

declare module "@mui/material/Paper" {
  interface PaperPropsVariantOverrides {
    centered: true;
  }
}

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
    MuiTableRow: {
      styleOverrides: {
        root: {
          "&.Mui-selected, &.Mui-selected:hover": {
            backgroundColor: theme.palette.grey[700],
          },
        },
      },
    },
    MuiPaper: {
      variants: [
        {
          props: { variant: "centered" },
          style: {
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            padding: theme.spacing(4),
          },
        },
      ],
    },
  },
});

export default theme;
