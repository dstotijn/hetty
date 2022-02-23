import ListIcon from "@mui/icons-material/List";
import { Button, Typography } from "@mui/material";
import Link from "next/link";
import React from "react";

import { Layout, Page } from "features/Layout";

function Index(): JSX.Element {
  return (
    <Layout page={Page.ProxySetup} title="Proxy setup">
      <Typography paragraph>Coming soon…</Typography>
      <Link href="/proxy/logs" passHref>
        <Button variant="contained" color="primary" component="a" size="large" startIcon={<ListIcon />}>
          View logs
        </Button>
      </Link>
    </Layout>
  );
}

export default Index;
