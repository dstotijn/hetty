import React from "react";
import { Button, Typography } from "@material-ui/core";
import ListIcon from "@material-ui/icons/List";
import Link from "next/link";

import Layout, { Page } from "../../components/Layout";

function Index(): JSX.Element {
  return (
    <Layout page={Page.ProxySetup} title="Proxy setup">
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
