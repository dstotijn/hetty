import { Box, Paper, Typography } from "@mui/material";
import { useRouter } from "next/router";

import { useInterceptedRequests } from "lib/InterceptedRequestsContext";
import RequestsTable from "lib/components/RequestsTable";

function Requests(): JSX.Element {
  const interceptedRequests = useInterceptedRequests();

  const router = useRouter();
  const activeId = router.query.id as string | undefined;

  const handleRowClick = (id: string) => {
    router.push(`/proxy/intercept?id=${id}`);
  };

  return (
    <Box>
      {interceptedRequests && interceptedRequests.length > 0 && (
        <RequestsTable requests={interceptedRequests} onRowClick={handleRowClick} activeRowId={activeId} />
      )}
      <Box sx={{ mt: 2, height: "100%" }}>
        {interceptedRequests?.length === 0 && (
          <Paper variant="centered">
            <Typography>No pending intercepted requests.</Typography>
          </Paper>
        )}
      </Box>
    </Box>
  );
}

export default Requests;
