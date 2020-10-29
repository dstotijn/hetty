import { Box } from "@material-ui/core";

import LogsOverview from "../../../components/reqlog/LogsOverview";
import Layout, { Page } from "../../../components/Layout";
import Search from "../../../components/reqlog/Search";

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
