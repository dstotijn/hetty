import LogsOverview from "../../../components/reqlog/LogsOverview";
import Layout, { Page } from "../../../components/Layout";

function ProxyLogs(): JSX.Element {
  return (
    <Layout page={Page.ProxyLogs} title="Proxy logs">
      <LogsOverview />
    </Layout>
  );
}

export default ProxyLogs;
