import {
  Alert,
  Box,
  BoxProps,
  Button,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  TextField,
  Typography,
} from "@mui/material";
import { ComponentType, FormEventHandler, useEffect, useRef, useState } from "react";
import { AllotmentProps, PaneProps } from "allotment/dist/types/src/allotment";

import { KeyValuePair, sortKeyValuePairs } from "./KeyValuePair";
import {
  GetSenderRequestQuery,
  HttpProtocol,
  useCreateOrUpdateSenderRequestMutation,
  useGetSenderRequestQuery,
  useSendRequestMutation,
} from "../../generated/graphql";
import EditRequestTabs from "./EditRequestTabs";
import Response from "./Response";

import "allotment/dist/style.css";
import { useRouter } from "next/router";

enum HttpMethod {
  Get = "GET",
  Post = "POST",
  Put = "PUT",
  Patch = "PATCH",
  Delete = "DELETE",
  Head = "HEAD",
  Options = "OPTIONS",
  Connect = "CONNECT",
  Trace = "TRACE",
}

enum HttpProto {
  Http1 = "HTTP/1.1",
  Http2 = "HTTP/2.0",
}

const httpProtoMap = new Map([
  [HttpProto.Http1, HttpProtocol.Http1],
  [HttpProto.Http2, HttpProtocol.Http2],
]);

function updateKeyPairItem(key: string, value: string, idx: number, items: any[]): any[] {
  const updated = [...items];
  updated[idx] = { key, value };

  // Append an empty key-value pair if the last item in the array isn't blank
  // anymore.
  if (items.length - 1 === idx && items[idx].key === "" && items[idx].value === "") {
    updated.push({ key: "", value: "" });
  }

  return updated;
}

function updateURLQueryParams(url: string, queryParams: KeyValuePair[]) {
  // Note: We don't use the `URL` interface, because we're potentially dealing
  // with malformed/incorrect URLs, which would yield TypeErrors when constructed
  // via `URL`.
  let newURL = url;

  const questionMarkIndex = url.indexOf("?");
  if (questionMarkIndex !== -1) {
    newURL = newURL.slice(0, questionMarkIndex);
  }

  const searchParams = new URLSearchParams();
  for (const { key, value } of queryParams.filter(({ key }) => key !== "")) {
    searchParams.append(key, value);
  }

  const rawQueryParams = decodeURI(searchParams.toString());

  if (rawQueryParams == "") {
    return newURL;
  }

  return newURL + "?" + rawQueryParams;
}

function queryParamsFromURL(url: string): KeyValuePair[] {
  const questionMarkIndex = url.indexOf("?");
  if (questionMarkIndex === -1) {
    return [];
  }

  const queryParams: KeyValuePair[] = [];

  const searchParams = new URLSearchParams(url.slice(questionMarkIndex + 1));
  for (let [key, value] of searchParams) {
    queryParams.push({ key, value });
  }

  return queryParams;
}

function EditRequest(): JSX.Element {
  const router = useRouter();
  const reqId = router.query.id as string | undefined;

  const [method, setMethod] = useState(HttpMethod.Get);
  const [url, setURL] = useState("");
  const [proto, setProto] = useState(HttpProto.Http2);
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
      console.log(senderRequest.response);
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

  const handleFormSubmit: FormEventHandler = (e) => {
    e.preventDefault();
    createOrUpdateRequestAndSend();
  };

  const isMountedRef = useRef(false);
  const [Allotment, setAllotment] = useState<
    (ComponentType<AllotmentProps> & { Pane: ComponentType<PaneProps> }) | null
  >(null);
  useEffect(() => {
    isMountedRef.current = true;
    import("allotment")
      .then((mod) => {
        if (!isMountedRef.current) {
          return;
        }
        setAllotment(mod.Allotment);
      })
      .catch((err) => console.error(err, `could not import allotment ${err.message}`));
    return () => {
      isMountedRef.current = false;
    };
  }, []);
  if (!Allotment) {
    return <div>Loading...</div>;
  }

  return (
    <Box display="flex" flexDirection="column" height="100%" gap={2}>
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

      <Box flex="1 auto" overflow="hidden">
        <Allotment>
          <Box pr={2} pb={2} height="100%" overflow="hidden">
            <Box height="100%" position="relative">
              <Typography variant="overline" color="textSecondary" sx={{ position: "absolute", right: 0, mt: 1.2 }}>
                Request
              </Typography>
              <EditRequestTabs
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
          </Box>
          <Box pb={2} pl={2} height="100%" overflow="hidden">
            <Box height="100%" position="relative">
              <Response response={response} />
            </Box>
          </Box>
        </Allotment>
      </Box>
    </Box>
  );
}

interface UrlBarProps extends BoxProps {
  method: HttpMethod;
  onMethodChange: (method: HttpMethod) => void;
  url: string;
  onUrlChange: (url: string) => void;
  proto: HttpProto;
  onProtoChange: (proto: HttpProto) => void;
}

function UrlBar(props: UrlBarProps) {
  const { method, onMethodChange, url, onUrlChange, proto, onProtoChange, ...other } = props;

  return (
    <Box {...other} sx={{ ...other.sx, display: "flex" }}>
      <FormControl>
        <InputLabel id="req-method-label">Method</InputLabel>
        <Select
          labelId="req-method-label"
          id="req-method"
          value={method}
          label="Method"
          onChange={(e) => onMethodChange(e.target.value as HttpMethod)}
          sx={{
            width: "8rem",
            ".MuiOutlinedInput-notchedOutline": {
              borderRightWidth: 0,
              borderTopRightRadius: 0,
              borderBottomRightRadius: 0,
            },
            "&:hover .MuiOutlinedInput-notchedOutline": {
              borderRightWidth: 1,
            },
          }}
        >
          {Object.values(HttpMethod).map((method) => (
            <MenuItem key={method} value={method}>
              {method}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
      <TextField
        label="URL"
        placeholder="E.g. “https://example.com/foobar”"
        value={url}
        onChange={(e) => onUrlChange(e.target.value)}
        required
        variant="outlined"
        InputLabelProps={{
          shrink: true,
        }}
        InputProps={{
          sx: {
            ".MuiOutlinedInput-notchedOutline": {
              borderRadius: 0,
            },
          },
        }}
        sx={{ flexGrow: 1 }}
      />
      <FormControl>
        <InputLabel id="req-proto-label">Protocol</InputLabel>
        <Select
          labelId="req-proto-label"
          id="req-proto"
          value={proto}
          label="Protocol"
          onChange={(e) => onProtoChange(e.target.value as HttpProto)}
          sx={{
            ".MuiOutlinedInput-notchedOutline": {
              borderLeftWidth: 0,
              borderTopLeftRadius: 0,
              borderBottomLeftRadius: 0,
            },
            "&:hover .MuiOutlinedInput-notchedOutline": {
              borderLeftWidth: 1,
            },
          }}
        >
          {Object.values(HttpProto).map((proto) => (
            <MenuItem key={proto} value={proto}>
              {proto}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </Box>
  );
}

export default EditRequest;
