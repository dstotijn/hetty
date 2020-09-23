import { Box, Typography } from "@material-ui/core";

import Layout from "../../components/Layout";

function Index(): JSX.Element {
  return (
    <Layout>
      <Box mb={2}>
        <Typography variant="h5">Sender</Typography>
      </Box>
      <Typography paragraph>Coming soon…</Typography>
    </Layout>
  );
}

export default Index;
