import { Layout, Page } from "features/Layout";
import RequestLogs from "features/reqlog";

function ProxyLogs(): JSX.Element {
  return (
    <Layout page={Page.ProxyLogs} title="Proxy logs">
      <RequestLogs />
    </Layout>
  );
}

export default ProxyLogs;
