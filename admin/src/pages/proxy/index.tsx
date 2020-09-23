import React from "react";
import { Box, Button, Typography } from "@material-ui/core";
import ListIcon from "@material-ui/icons/List";
import Link from "next/link";

import Layout from "../../components/Layout";

function Index(): JSX.Element {
  return (
    <Layout>
      <Box mb={2}>
        <Typography variant="h5">Proxy setup</Typography>
      </Box>
      <Typography paragraph>Coming soon…</Typography>
      <Link href="/proxy/logs" passHref>
        <Button
          variant="contained"
          color="secondary"
          component="a"
          size="large"
          startIcon={<ListIcon />}
        >
          View logs
        </Button>
      </Link>
    </Layout>
  );
}

export default Index;
