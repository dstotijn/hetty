import { Layout, Page } from "features/Layout";
import Sender from "features/sender";

function Index(): JSX.Element {
  return (
    <Layout page={Page.Sender} title="Sender">
      <Sender />
    </Layout>
  );
}

export default Index;
