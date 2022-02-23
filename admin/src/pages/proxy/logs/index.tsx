import { Box } from "@mui/material";

import { Layout, Page } from "features/Layout";
import LogsOverview from "features/reqlog/components/LogsOverview";
import Search from "features/reqlog/components/Search";

function ProxyLogs(): JSX.Element {
  return (
    <Layout page={Page.ProxyLogs} title="Proxy logs">
      <Box mb={2}>
        <Search />
      </Box>
      <LogsOverview />
    </Layout>
  );
}

export default ProxyLogs;
