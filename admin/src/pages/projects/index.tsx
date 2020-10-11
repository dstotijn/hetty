import { Box, Divider, Grid, Typography } from "@material-ui/core";
import Layout, { Page } from "../../components/Layout";
import NewProject from "../../components/projects/NewProject";
import ProjectList from "../../components/projects/ProjectList";

function Index(): JSX.Element {
  return (
    <Layout page={Page.Projects} title="Projects">
      <Box p={4}>
        <Box mb={3}>
          <Typography variant="h4">Projects</Typography>
        </Box>
        <Typography paragraph>
          Projects contain settings and data generated/processed by Hetty. They
          are stored as SQLite database files on disk.
        </Typography>
        <Box my={4}>
          <Divider />
        </Box>
        <Box mb={8}>
          <NewProject />
        </Box>
        <Grid container>
          <Grid item xs={12} sm={8} md={6} lg={6}>
            <ProjectList />
          </Grid>
        </Grid>
      </Box>
    </Layout>
  );
}

export default Index;
