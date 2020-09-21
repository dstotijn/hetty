import { useState } from "react";
import { Box } from "@material-ui/core";

import RequestList from "../../components/reqlog/RequestList";
import LogDetail from "../../components/reqlog/LogDetail";
import LogsOverview from "../../components/reqlog/LogsOverview";
import Layout from "../../components/Layout";

function Logs(): JSX.Element {
  return (
    <Layout>
      <LogsOverview />
    </Layout>
  );
}

export default Logs;
