import CancelIcon from "@mui/icons-material/Cancel";
import SendIcon from "@mui/icons-material/Send";
import { Alert, Box, Button, CircularProgress, Typography } from "@mui/material";
import { useRouter } from "next/router";
import React, { useEffect, useState } from "react";

import { useInterceptedRequests } from "lib/InterceptedRequestsContext";
import { KeyValuePair, sortKeyValuePairs } from "lib/components/KeyValuePair";
import RequestTabs from "lib/components/RequestTabs";
import Response from "lib/components/Response";
import SplitPane from "lib/components/SplitPane";
import UrlBar, { HttpMethod, HttpProto, httpProtoMap } from "lib/components/UrlBar";
import {
  HttpProtocol,
  HttpRequest,
  useCancelRequestMutation,
  useGetInterceptedRequestQuery,
  useModifyRequestMutation,
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
  const [headers, setHeaders] = useState<KeyValuePair[]>([{ key: "", value: "" }]);
  const [body, setBody] = useState("");

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

  const handleHeaderChange = (key: string, value: string, idx: number) => {
    setHeaders((prev) => updateKeyPairItem(key, value, idx, prev));
  };
  const handleHeaderDelete = (idx: number) => {
    setHeaders((prev) => prev.slice(0, idx).concat(prev.slice(idx + 1, prev.length)));
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
      setBody(interceptedRequest.body || "");

      const newQueryParams = queryParamsFromURL(interceptedRequest.url);
      // Push empty row.
      newQueryParams.push({ key: "", value: "" });
      setQueryParams(newQueryParams);

      const newHeaders = sortKeyValuePairs(interceptedRequest.headers || []);
      setHeaders([...newHeaders.map(({ key, value }) => ({ key, value })), { key: "", value: "" }]);
    },
  });
  const interceptedReq = reqId ? getReqResult?.data?.interceptedRequest : undefined;

  const [modifyRequest, modifyResult] = useModifyRequestMutation();
  const [cancelRequest, cancelResult] = useCancelRequestMutation();

  const onActionCompleted = () => {
    setURL("");
    setMethod(HttpMethod.Get);
    setBody("");
    setQueryParams([]);
    setHeaders([]);
    router.replace(`/proxy/intercept`);
  };

  const handleFormSubmit: React.FormEventHandler = (e) => {
    e.preventDefault();

    if (!interceptedReq) {
      return;
    }

    modifyRequest({
      variables: {
        request: {
          id: interceptedReq.id,
          url,
          method,
          proto: httpProtoMap.get(proto) || HttpProtocol.Http20,
          headers: headers.filter((kv) => kv.key !== ""),
          body: body || undefined,
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
  };

  const handleCancelClick = () => {
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
          <Button
            variant="contained"
            disableElevation
            type="submit"
            disabled={!interceptedReq || modifyResult.loading || cancelResult.loading}
            startIcon={modifyResult.loading ? <CircularProgress size={22} /> : <SendIcon />}
          >
            Send
          </Button>
          <Button
            variant="contained"
            color="error"
            disableElevation
            onClick={handleCancelClick}
            disabled={!interceptedReq || modifyResult.loading || cancelResult.loading}
            startIcon={cancelResult.loading ? <CircularProgress size={22} /> : <CancelIcon />}
          >
            Cancel
          </Button>
        </Box>
        {modifyResult.error && (
          <Alert severity="error" sx={{ mt: 1 }}>
            {modifyResult.error.message}
          </Alert>
        )}
        {cancelResult.error && (
          <Alert severity="error" sx={{ mt: 1 }}>
            {cancelResult.error.message}
          </Alert>
        )}
      </Box>

      <Box flex="1 auto" position="relative">
        <SplitPane split="vertical" size={"50%"}>
          <Box sx={{ height: "100%", mr: 2, pb: 2, position: "relative" }}>
            <Typography variant="overline" color="textSecondary" sx={{ position: "absolute", right: 0, mt: 1.2 }}>
              Request
            </Typography>
            <RequestTabs
              queryParams={interceptedReq ? queryParams : []}
              headers={interceptedReq ? headers : []}
              body={body}
              onQueryParamChange={interceptedReq ? handleQueryParamChange : undefined}
              onQueryParamDelete={interceptedReq ? handleQueryParamDelete : undefined}
              onHeaderChange={interceptedReq ? handleHeaderChange : undefined}
              onHeaderDelete={interceptedReq ? handleHeaderDelete : undefined}
              onBodyChange={interceptedReq ? setBody : undefined}
            />
          </Box>
          <Box sx={{ height: "100%", position: "relative", ml: 2, pb: 2 }}>
            <Response response={null} />
          </Box>
        </SplitPane>
      </Box>
    </Box>
  );
}

export default EditRequest;
