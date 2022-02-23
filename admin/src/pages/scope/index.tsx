import { Box, Divider, Grid, Typography } from "@mui/material";
import React from "react";

import { Layout, Page } from "features/Layout";
import AddRule from "features/scope/components/AddRule";
import Rules from "features/scope/components/Rules";

function Index(): JSX.Element {
  return (
    <Layout page={Page.Scope} title="Scope">
      <Box p={4}>
        <Box mb={3}>
          <Typography variant="h4">Scope</Typography>
        </Box>
        <Typography paragraph>
          Scope rules are used by various modules in Hetty and can influence their behavior. For example: the Proxy logs
          module can match incoming requests against scope rules and decide its behavior (e.g. log or bypass) based on
          the outcome of the match. All scope configuration is stored per project.
        </Typography>
        <Box my={4}>
          <Divider />
        </Box>
        <Grid container>
          <Grid item xs={12} sm={12} md={8} lg={6}>
            <AddRule />
            <Box my={4}>
              <Divider />
            </Box>
            <Rules />
          </Grid>
        </Grid>
      </Box>
    </Layout>
  );
}

export default Index;
