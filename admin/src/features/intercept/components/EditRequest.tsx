import CancelIcon from "@mui/icons-material/Cancel";
import DownloadIcon from "@mui/icons-material/Download";
import SendIcon from "@mui/icons-material/Send";
import SettingsIcon from "@mui/icons-material/Settings";
import { Alert, Box, Button, CircularProgress, IconButton, Tooltip, Typography } from "@mui/material";
import { useRouter } from "next/router";
import React, { useEffect, useState } from "react";

import { useInterceptedRequests } from "lib/InterceptedRequestsContext";
import { KeyValuePair } from "lib/components/KeyValuePair";
import Link from "lib/components/Link";
import RequestTabs from "lib/components/RequestTabs";
import ResponseStatus from "lib/components/ResponseStatus";
import ResponseTabs from "lib/components/ResponseTabs";
import UrlBar, { HttpMethod, HttpProto, httpProtoMap } from "lib/components/UrlBar";
import {
  HttpProtocol,
  HttpRequest,
  useCancelRequestMutation,
  useCancelResponseMutation,
  useGetInterceptedRequestQuery,
  useModifyRequestMutation,
  useModifyResponseMutation,
} from "lib/graphql/generated";
import { queryParamsFromURL } from "lib/queryParamsFromURL";
import updateKeyPairItem from "lib/updateKeyPairItem";
import updateURLQueryParams from "lib/updateURLQueryParams";

function EditRequest(): JSX.Element {
  const router = useRouter();
  const interceptedRequests = useInterceptedRequests();

  useEffect(() => {
    // If there's no request selected and there are pending reqs, navigate to
    // the first one in the list. This helps you quickly review/handle reqs
    // without having to manually select the next one in the requests table.
    if (router.isReady && !router.query.id && interceptedRequests?.length) {
      const req = interceptedRequests[0];
      router.replace(`/proxy/intercept?id=${req.id}`);
    }
  }, [router, interceptedRequests]);

  const reqId = router.query.id as string | undefined;

  const [method, setMethod] = useState(HttpMethod.Get);
  const [url, setURL] = useState("");
  const [proto, setProto] = useState(HttpProto.Http20);
  const [queryParams, setQueryParams] = useState<KeyValuePair[]>([{ key: "", value: "" }]);
  const [reqHeaders, setReqHeaders] = useState<KeyValuePair[]>([{ key: "", value: "" }]);
  const [resHeaders, setResHeaders] = useState<KeyValuePair[]>([{ key: "", value: "" }]);
  const [reqBody, setReqBody] = useState("");
  const [resBody, setResBody] = useState("");

  const handleQueryParamChange = (key: string, value: string, idx: number) => {
    setQueryParams((prev) => {
      const updated = updateKeyPairItem(key, value, idx, prev);
      setURL((prev) => updateURLQueryParams(prev, updated));
      return updated;
    });
  };
  const handleQueryParamDelete = (idx: number) => {
    setQueryParams((prev) => {
      const updated = prev.slice(0, idx).concat(prev.slice(idx + 1, prev.length));
      setURL((prev) => updateURLQueryParams(prev, updated));
      return updated;
    });
  };

  const handleReqHeaderChange = (key: string, value: string, idx: number) => {
    setReqHeaders((prev) => updateKeyPairItem(key, value, idx, prev));
  };
  const handleReqHeaderDelete = (idx: number) => {
    setReqHeaders((prev) => prev.slice(0, idx).concat(prev.slice(idx + 1, prev.length)));
  };

  const handleResHeaderChange = (key: string, value: string, idx: number) => {
    setResHeaders((prev) => updateKeyPairItem(key, value, idx, prev));
  };
  const handleResHeaderDelete = (idx: number) => {
    setResHeaders((prev) => prev.slice(0, idx).concat(prev.slice(idx + 1, prev.length)));
  };

  const handleURLChange = (url: string) => {
    setURL(url);

    const questionMarkIndex = url.indexOf("?");
    if (questionMarkIndex === -1) {
      setQueryParams([{ key: "", value: "" }]);
      return;
    }

    const newQueryParams = queryParamsFromURL(url);
    // Push empty row.
    newQueryParams.push({ key: "", value: "" });
    setQueryParams(newQueryParams);
  };

  const getReqResult = useGetInterceptedRequestQuery({
    variables: { id: reqId as string },
    skip: reqId === undefined,
    onCompleted: ({ interceptedRequest }) => {
      if (!interceptedRequest) {
        return;
      }

      setURL(interceptedRequest.url);
      setMethod(interceptedRequest.method);
      setReqBody(interceptedRequest.body || "");

      const newQueryParams = queryParamsFromURL(interceptedRequest.url);
      // Push empty row.
      newQueryParams.push({ key: "", value: "" });
      setQueryParams(newQueryParams);

      const newReqHeaders = interceptedRequest.headers || [];
      setReqHeaders([...newReqHeaders.map(({ key, value }) => ({ key, value })), { key: "", value: "" }]);

      setResBody(interceptedRequest.response?.body || "");
      const newResHeaders = interceptedRequest.response?.headers || [];
      setResHeaders([...newResHeaders.map(({ key, value }) => ({ key, value })), { key: "", value: "" }]);
    },
  });
  const interceptedReq =
    reqId && !getReqResult?.data?.interceptedRequest?.response ? getReqResult?.data?.interceptedRequest : undefined;
  const interceptedRes = reqId ? getReqResult?.data?.interceptedRequest?.response : undefined;

  const [modifyRequest, modifyReqResult] = useModifyRequestMutation();
  const [cancelRequest, cancelReqResult] = useCancelRequestMutation();

  const [modifyResponse, modifyResResult] = useModifyResponseMutation();
  const [cancelResponse, cancelResResult] = useCancelResponseMutation();

  const onActionCompleted = () => {
    setURL("");
    setMethod(HttpMethod.Get);
    setReqBody("");
    setQueryParams([]);
    setReqHeaders([]);
    router.replace(`/proxy/intercept`);
  };

  const handleFormSubmit: React.FormEventHandler = (e) => {
    e.preventDefault();

    if (interceptedReq) {
      modifyRequest({
        variables: {
          request: {
            id: interceptedReq.id,
            url,
            method,
            proto: httpProtoMap.get(proto) || HttpProtocol.Http20,
            headers: reqHeaders.filter((kv) => kv.key !== ""),
            body: reqBody || undefined,
          },
        },
        update(cache) {
          cache.modify({
            fields: {
              interceptedRequests(existing: HttpRequest[], { readField }) {
                return existing.filter((ref) => interceptedReq.id !== readField("id", ref));
              },
            },
          });
        },
        onCompleted: onActionCompleted,
      });
    }

    if (interceptedRes) {
      modifyResponse({
        variables: {
          response: {
            requestID: interceptedRes.id,
            proto: interceptedRes.proto, // TODO: Allow modifying
            statusCode: interceptedRes.statusCode, // TODO: Allow modifying
            statusReason: interceptedRes.statusReason, // TODO: Allow modifying
            headers: resHeaders.filter((kv) => kv.key !== ""),
            body: resBody || undefined,
          },
        },
        update(cache) {
          cache.modify({
            fields: {
              interceptedRequests(existing: HttpRequest[], { readField }) {
                return existing.filter((ref) => interceptedRes.id !== readField("id", ref));
              },
            },
          });
        },
        onCompleted: onActionCompleted,
      });
    }
  };

  const handleReqCancelClick = () => {
    if (!interceptedReq) {
      return;
    }

    cancelRequest({
      variables: {
        id: interceptedReq.id,
      },
      update(cache) {
        cache.modify({
          fields: {
            interceptedRequests(existing: HttpRequest[], { readField }) {
              return existing.filter((ref) => interceptedReq.id !== readField("id", ref));
            },
          },
        });
      },
      onCompleted: onActionCompleted,
    });
  };

  const handleResCancelClick = () => {
    if (!interceptedRes) {
      return;
    }

    cancelResponse({
      variables: {
        requestID: interceptedRes.id,
      },
      update(cache) {
        cache.modify({
          fields: {
            interceptedRequests(existing: HttpRequest[], { readField }) {
              return existing.filter((ref) => interceptedRes.id !== readField("id", ref));
            },
          },
        });
      },
      onCompleted: onActionCompleted,
    });
  };

  return (
    <Box display="flex" flexDirection="column" height="100%" gap={2}>
      <Box component="form" autoComplete="off" onSubmit={handleFormSubmit}>
        <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
          <UrlBar
            method={method}
            onMethodChange={interceptedReq ? setMethod : undefined}
            url={url.toString()}
            onUrlChange={interceptedReq ? handleURLChange : undefined}
            proto={proto}
            onProtoChange={interceptedReq ? setProto : undefined}
            sx={{ flex: "1 auto" }}
          />
          {!interceptedRes && (
            <>
              <Button
                variant="contained"
                disableElevation
                type="submit"
                disabled={!interceptedReq || modifyReqResult.loading || cancelReqResult.loading}
                startIcon={modifyReqResult.loading ? <CircularProgress size={22} /> : <SendIcon />}
              >
                Send
              </Button>
              <Button
                variant="contained"
                color="error"
                disableElevation
                onClick={handleReqCancelClick}
                disabled={!interceptedReq || modifyReqResult.loading || cancelReqResult.loading}
                startIcon={cancelReqResult.loading ? <CircularProgress size={22} /> : <CancelIcon />}
              >
                Cancel
              </Button>
            </>
          )}
          {interceptedRes && (
            <>
              <Button
                variant="contained"
                disableElevation
                type="submit"
                disabled={modifyResResult.loading || cancelResResult.loading}
                endIcon={modifyResResult.loading ? <CircularProgress size={22} /> : <DownloadIcon />}
              >
                Receive
              </Button>
              <Button
                variant="contained"
                color="error"
                disableElevation
                onClick={handleResCancelClick}
                disabled={modifyResResult.loading || cancelResResult.loading}
                endIcon={cancelResResult.loading ? <CircularProgress size={22} /> : <CancelIcon />}
              >
                Cancel
              </Button>
            </>
          )}
          <Tooltip title="Intercept settings">
            <IconButton LinkComponent={Link} href="/settings#intercept">
              <SettingsIcon />
            </IconButton>
          </Tooltip>
        </Box>
        {modifyReqResult.error && (
          <Alert severity="error" sx={{ mt: 1 }}>
            {modifyReqResult.error.message}
          </Alert>
        )}
        {cancelReqResult.error && (
          <Alert severity="error" sx={{ mt: 1 }}>
            {cancelReqResult.error.message}
          </Alert>
        )}
      </Box>

      <Box flex="1 auto" overflow="scroll">
        {interceptedReq && (
          <Box sx={{ height: "100%", pb: 2 }}>
            <Typography variant="overline" color="textSecondary" sx={{ position: "absolute", right: 0, mt: 1.2 }}>
              Request
            </Typography>
            <RequestTabs
              queryParams={interceptedReq ? queryParams : []}
              headers={interceptedReq ? reqHeaders : []}
              body={reqBody}
              onQueryParamChange={interceptedReq ? handleQueryParamChange : undefined}
              onQueryParamDelete={interceptedReq ? handleQueryParamDelete : undefined}
              onHeaderChange={interceptedReq ? handleReqHeaderChange : undefined}
              onHeaderDelete={interceptedReq ? handleReqHeaderDelete : undefined}
              onBodyChange={interceptedReq ? setReqBody : undefined}
            />
          </Box>
        )}
        {interceptedRes && (
          <Box sx={{ height: "100%", pb: 2 }}>
            <Box sx={{ position: "absolute", right: 0, mt: 1.4 }}>
              <Typography variant="overline" color="textSecondary" sx={{ float: "right", ml: 3 }}>
                Response
              </Typography>
              {interceptedRes && (
                <Box sx={{ float: "right", mt: 0.2 }}>
                  <ResponseStatus
                    proto={interceptedRes.proto}
                    statusCode={interceptedRes.statusCode}
                    statusReason={interceptedRes.statusReason}
                  />
                </Box>
              )}
            </Box>
            <ResponseTabs
              headers={interceptedRes ? resHeaders : []}
              body={resBody}
              onHeaderChange={interceptedRes ? handleResHeaderChange : undefined}
              onHeaderDelete={interceptedRes ? handleResHeaderDelete : undefined}
              onBodyChange={interceptedRes ? setResBody : undefined}
              hasResponse={interceptedRes !== undefined && interceptedRes !== null}
            />
          </Box>
        )}
      </Box>
    </Box>
  );
}

export default EditRequest;
