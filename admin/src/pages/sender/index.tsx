import { Box, Typography } from "@material-ui/core";

import Layout, { Page } from "../../components/Layout";

function Index(): JSX.Element {
  return (
    <Layout page={Page.Sender} title="Sender">
      <Typography paragraph>Coming soon…</Typography>
    </Layout>
  );
}

export default Index;
