import { Typography, Box } from "@material-ui/core";

import LogsOverview from "../../../components/reqlog/LogsOverview";
import Layout from "../../../components/Layout";

function ProxyLogs(): JSX.Element {
  return (
    <Layout>
      <Box mb={2}>
        <Typography variant="h5">Proxy logs</Typography>
      </Box>
      <LogsOverview />
    </Layout>
  );
}

export default ProxyLogs;
