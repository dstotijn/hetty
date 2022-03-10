import AltRouteIcon from "@mui/icons-material/AltRoute";
import DeleteIcon from "@mui/icons-material/Delete";
import { Alert } from "@mui/lab";
import { Badge, Button, IconButton, Tooltip } from "@mui/material";
import Link from "next/link";

import { useInterceptedRequests } from "lib/InterceptedRequestsContext";
import { ConfirmationDialog, useConfirmationDialog } from "lib/components/ConfirmationDialog";
import { HttpRequestLogsDocument, useClearHttpRequestLogMutation } from "lib/graphql/generated";

function Actions(): JSX.Element {
  const interceptedRequests = useInterceptedRequests();
  const [clearHTTPRequestLog, clearLogsResult] = useClearHttpRequestLogMutation({
    refetchQueries: [{ query: HttpRequestLogsDocument }],
  });
  const clearHTTPConfirmationDialog = useConfirmationDialog();

  return (
    <div>
      <ConfirmationDialog
        isOpen={clearHTTPConfirmationDialog.isOpen}
        onClose={clearHTTPConfirmationDialog.close}
        onConfirm={clearHTTPRequestLog}
      >
        All proxy logs are going to be removed. This action cannot be undone.
      </ConfirmationDialog>

      {clearLogsResult.error && <Alert severity="error">Failed to clear HTTP logs: {clearLogsResult.error}</Alert>}

      <Link href="/proxy/intercept/?id=" passHref>
        <Button
          variant="contained"
          disabled={interceptedRequests === null || interceptedRequests.length === 0}
          color="primary"
          component="a"
          size="large"
          startIcon={
            <Badge color="error" badgeContent={interceptedRequests?.length || 0}>
              <AltRouteIcon />
            </Badge>
          }
          sx={{ mr: 1 }}
        >
          Review Intercepted…
        </Button>
      </Link>

      <Tooltip title="Clear all">
        <IconButton onClick={clearHTTPConfirmationDialog.open}>
          <DeleteIcon />
        </IconButton>
      </Tooltip>
    </div>
  );
}

export default Actions;
