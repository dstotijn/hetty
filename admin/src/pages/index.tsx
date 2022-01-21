import {
  Box,
  Button,
  createStyles,
  makeStyles,
  Theme,
  Typography,
} from "@material-ui/core";
import FolderIcon from "@material-ui/icons/Folder";
import Link from "next/link";

import { useRouter } from "next/router";

import Layout, { Page } from "../components/Layout";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    titleHighlight: {
      color: theme.palette.secondary.main,
    },
    subtitle: {
      fontSize: "1.6rem",
      width: "60%",
      lineHeight: 2,
      marginBottom: theme.spacing(5),
    },
    button: {
      marginRight: theme.spacing(2),
    },
  })
);

function Index(): JSX.Element {
  const classes = useStyles();

  return (
    <Layout page={Page.Home} title="">
      <Box p={4}>
        <Box mb={4} width="60%">
          <Typography variant="h2">
            <span className={classes.titleHighlight}>Hetty://</span>
            <br />
            The simple HTTP toolkit for security research.
          </Typography>
        </Box>

        <Typography className={classes.subtitle} paragraph>
          What if security testing was intuitive, powerful, and good looking?
          What if it was <strong>free</strong>, instead of $400 per year?{" "}
          <span className={classes.titleHighlight}>Hetty</span> is listening on{" "}
          <code>:8080</code>…
        </Typography>

        <Link href="/projects" passHref>
          <Button
            className={classes.button}
            variant="contained"
            color="secondary"
            component="a"
            size="large"
            startIcon={<FolderIcon />}
          >
            Manage projects
          </Button>
        </Link>
      </Box>
    </Layout>
  );
}

export default Index;
