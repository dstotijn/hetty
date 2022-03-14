import { Layout, Page } from "features/Layout";
import Settings from "features/settings/components/Settings";

function Index(): JSX.Element {
  return (
    <Layout page={Page.Settings} title="Settings">
      <Settings />
    </Layout>
  );
}

export default Index;
