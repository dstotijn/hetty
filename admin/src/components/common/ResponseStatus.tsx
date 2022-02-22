import { Typography } from "@mui/material";

import HttpStatusIcon from "./HttpStatusIcon";
import { HttpProtocol } from "../../generated/graphql";

type ResponseStatusProps = {
  proto: HttpProtocol;
  statusCode: number;
  statusReason: string;
};

function mapProto(proto: HttpProtocol): string {
  switch (proto) {
    case HttpProtocol.Http1:
      return "HTTP/1.1";
    case HttpProtocol.Http2:
      return "HTTP/2.0";
    default:
      return proto;
  }
}

function ResponseStatus({ proto, statusCode, statusReason }: ResponseStatusProps): JSX.Element {
  return (
    <Typography variant="h6" style={{ fontSize: "1rem", whiteSpace: "nowrap" }}>
      <HttpStatusIcon status={statusCode} />{" "}
      <Typography component="span" color="textSecondary">
        <Typography component="span" color="textSecondary" style={{ fontFamily: "'JetBrains Mono', monospace" }}>
          {mapProto(proto)}
        </Typography>
      </Typography>{" "}
      {statusCode} {statusReason}
    </Typography>
  );
}

export default ResponseStatus;
