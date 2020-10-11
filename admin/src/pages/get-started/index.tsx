import { Box, Link as MaterialLink, Typography } from "@material-ui/core";
import Link from "next/link";

import React from "react";

import Layout, { Page } from "../../components/Layout";

function Index(): JSX.Element {
  return (
    <Layout page={Page.GetStarted} title="Get started">
      <Box p={4}>
        <Box mb={3}>
          <Typography variant="h4">Get started</Typography>
        </Box>
        <Typography paragraph>
          You’ve loaded a (new) project. What’s next? You can now use the MITM
          proxy and review HTTP requests and responses via the{" "}
          <Link href="/proxy/logs" passHref>
            <MaterialLink color="secondary">Proxy logs</MaterialLink>
          </Link>
          . Stuck? Ask for help on the{" "}
          <MaterialLink
            href="https://github.com/dstotijn/hetty/discussions"
            color="secondary"
            target="_blank"
          >
            Discussions forum
          </MaterialLink>
          .
        </Typography>
      </Box>
    </Layout>
  );
}

export default Index;
