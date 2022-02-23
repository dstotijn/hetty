import { Box, Typography } from "@mui/material";

import { sortKeyValuePairs } from "./KeyValuePair";
import ResponseTabs from "./ResponseTabs";

import ResponseStatus from "lib/components/ResponseStatus";
import { HttpResponseLog } from "lib/graphql/generated";

interface ResponseProps {
  response?: HttpResponseLog | null;
}

function Response({ response }: ResponseProps): JSX.Element {
  return (
    <Box height="100%">
      <div>
        <Box sx={{ position: "absolute", right: 0, mt: 1.4 }}>
          <Typography variant="overline" color="textSecondary" sx={{ float: "right", ml: 3 }}>
            Response
          </Typography>
          {response && (
            <Box sx={{ float: "right", mt: 0.2 }}>
              <ResponseStatus
                proto={response.proto}
                statusCode={response.statusCode}
                statusReason={response.statusReason}
              />
            </Box>
          )}
        </Box>
      </div>
      <ResponseTabs
        body={response?.body}
        headers={sortKeyValuePairs(response?.headers || [])}
        hasResponse={response !== undefined && response !== null}
      />
    </Box>
  );
}

export default Response;
