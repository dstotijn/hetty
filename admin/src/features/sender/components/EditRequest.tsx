import AddIcon from "@mui/icons-material/Add";
import { Alert, Box, Button, Fab, Tooltip, Typography, useTheme } from "@mui/material";
import { useRouter } from "next/router";
import React, { useState } from "react";

import { KeyValuePair, sortKeyValuePairs } from "lib/components/KeyValuePair";
import RequestTabs from "lib/components/RequestTabs";
import Response from "lib/components/Response";
import SplitPane from "lib/components/SplitPane";
import UrlBar, { HttpMethod, HttpProto, httpProtoMap } from "lib/components/UrlBar";
import {
  GetSenderRequestQuery,
  useCreateOrUpdateSenderRequestMutation,
  useGetSenderRequestQuery,
  useSendRequestMutation,
} from "lib/graphql/generated";
import { queryParamsFromURL } from "lib/queryParamsFromURL";
import updateKeyPairItem from "lib/updateKeyPairItem";
import updateURLQueryParams from "lib/updateURLQueryParams";

const defaultMethod = HttpMethod.Get;
const defaultProto = HttpProto.Http20;
const emptyKeyPair = [{ key: "", value: "" }];

function EditRequest(): JSX.Element {
  const router = useRouter();
  const reqId = router.query.id as string | undefined;

  const theme = useTheme();

  const [method, setMethod] = useState(defaultMethod);
  const [url, setURL] = useState("");
  const [proto, setProto] = useState(defaultProto);
  const [queryParams, setQueryParams] = useState<KeyValuePair[]>(emptyKeyPair);
  const [headers, setHeaders] = useState<KeyValuePair[]>(emptyKeyPair);
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

  const [response, setResponse] = useState<NonNullable<GetSenderRequestQuery["senderRequest"]>["response"]>(null);
  const getReqResult = useGetSenderRequestQuery({
    variables: { id: reqId as string },
    skip: reqId === undefined,
    onCompleted: ({ senderRequest }) => {
      if (!senderRequest) {
        return;
      }

      setURL(senderRequest.url);
      setMethod(senderRequest.method);
      setBody(senderRequest.body || "");

      const newQueryParams = queryParamsFromURL(senderRequest.url);
      // Push empty row.
      newQueryParams.push({ key: "", value: "" });
      setQueryParams(newQueryParams);

      const newHeaders = sortKeyValuePairs(senderRequest.headers || []);
      setHeaders([...newHeaders.map(({ key, value }) => ({ key, value })), { key: "", value: "" }]);
      setResponse(senderRequest.response);
    },
  });

  const [createOrUpdateRequest, createResult] = useCreateOrUpdateSenderRequestMutation();
  const [sendRequest, sendResult] = useSendRequestMutation();

  const createOrUpdateRequestAndSend = () => {
    const senderReq = getReqResult?.data?.senderRequest;
    createOrUpdateRequest({
      variables: {
        request: {
          // Update existing sender request if it was cloned from a request log
          // and it doesn't have a response body yet (e.g. not sent yet).
          ...(senderReq && senderReq.sourceRequestLogID && !senderReq.response && { id: senderReq.id }),
          url,
          method,
          proto: httpProtoMap.get(proto),
          headers: headers.filter((kv) => kv.key !== ""),
          body: body || undefined,
        },
      },
      onCompleted: ({ createOrUpdateSenderRequest }) => {
        const { id } = createOrUpdateSenderRequest;
        sendRequestAndPushRoute(id);
      },
    });
  };

  const sendRequestAndPushRoute = (id: string) => {
    sendRequest({
      errorPolicy: "all",
      onCompleted: () => {
        router.push(`/sender?id=${id}`);
      },
      variables: {
        id,
      },
    });
  };

  const handleFormSubmit: React.FormEventHandler = (e) => {
    e.preventDefault();
    createOrUpdateRequestAndSend();
  };

  const handleNewRequest = () => {
    setURL("");
    setMethod(defaultMethod);
    setProto(defaultProto);
    setQueryParams(emptyKeyPair);
    setHeaders(emptyKeyPair);
    setBody("");
    setResponse(null);
    router.push(`/sender`);
  };

  return (
    <Box display="flex" flexDirection="column" height="100%" gap={2}>
      <Box sx={{ position: "absolute", bottom: theme.spacing(2), right: theme.spacing(2) }}>
        <Tooltip title="New request">
          <Fab color="primary" onClick={handleNewRequest}>
            <AddIcon />
          </Fab>
        </Tooltip>
      </Box>
      <Box component="form" autoComplete="off" onSubmit={handleFormSubmit}>
        <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
          <UrlBar
            method={method}
            onMethodChange={setMethod}
            url={url.toString()}
            onUrlChange={handleURLChange}
            proto={proto}
            onProtoChange={setProto}
            sx={{ flex: "1 auto" }}
          />
          <Button
            variant="contained"
            disableElevation
            sx={{ width: "8rem" }}
            type="submit"
            disabled={createResult.loading || sendResult.loading}
          >
            Send
          </Button>
        </Box>
        {createResult.error && (
          <Alert severity="error" sx={{ mt: 1 }}>
            {createResult.error.message}
          </Alert>
        )}
        {sendResult.error && (
          <Alert severity="error" sx={{ mt: 1 }}>
            {sendResult.error.message}
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
              queryParams={queryParams}
              headers={headers}
              body={body}
              onQueryParamChange={handleQueryParamChange}
              onQueryParamDelete={handleQueryParamDelete}
              onHeaderChange={handleHeaderChange}
              onHeaderDelete={handleHeaderDelete}
              onBodyChange={setBody}
            />
          </Box>
          <Box sx={{ height: "100%", position: "relative", ml: 2, pb: 2 }}>
            <Response response={response} />
          </Box>
        </SplitPane>
      </Box>
    </Box>
  );
}

export default EditRequest;
