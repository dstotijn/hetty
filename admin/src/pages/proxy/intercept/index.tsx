import { Layout, Page } from "features/Layout";
import Intercept from "features/intercept/components/Intercept";

function ProxyIntercept(): JSX.Element {
  return (
    <Layout page={Page.Intercept} title="Proxy intercept">
      <Intercept />
    </Layout>
  );
}

export default ProxyIntercept;
