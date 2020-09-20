import { useState } from "react";
import { Box } from "@material-ui/core";

import RequestList from "../../components/reqlog/RequestList";
import LogDetail from "../../components/reqlog/LogDetail";

function Logs(): JSX.Element {
  const [detailReqLogId, setDetailReqLogId] = useState<string>();

  const handleLogClick = (reqId: string) => setDetailReqLogId(reqId);

  return (
    <div>
      <Box minHeight="375px" maxHeight="33vh" overflow="scroll">
        <RequestList onLogClick={handleLogClick} />
      </Box>
      <Box>{detailReqLogId && <LogDetail requestId={detailReqLogId} />}</Box>
    </div>
  );
}

export default Logs;
