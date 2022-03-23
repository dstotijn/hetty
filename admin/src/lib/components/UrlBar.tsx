import { Box, BoxProps, FormControl, InputLabel, MenuItem, Select, TextField } from "@mui/material";

import { HttpProtocol } from "lib/graphql/generated";

export enum HttpMethod {
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

export enum HttpProto {
  Http10 = "HTTP/1.0",
  Http11 = "HTTP/1.1",
  Http20 = "HTTP/2.0",
}

export const httpProtoMap = new Map([
  [HttpProto.Http10, HttpProtocol.Http10],
  [HttpProto.Http11, HttpProtocol.Http11],
  [HttpProto.Http20, HttpProtocol.Http20],
]);

interface UrlBarProps extends BoxProps {
  method: HttpMethod;
  onMethodChange?: (method: HttpMethod) => void;
  url: string;
  onUrlChange?: (url: string) => void;
  proto: HttpProto;
  onProtoChange?: (proto: HttpProto) => void;
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
          disabled={!onMethodChange}
          onChange={(e) => onMethodChange && onMethodChange(e.target.value as HttpMethod)}
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
        disabled={!onUrlChange}
        onChange={(e) => onUrlChange && onUrlChange(e.target.value)}
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
          disabled={!onProtoChange}
          onChange={(e) => onProtoChange && onProtoChange(e.target.value as HttpProto)}
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

export default UrlBar;
