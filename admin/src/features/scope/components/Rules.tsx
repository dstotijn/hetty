import { Alert } from "@mui/lab";
import { CircularProgress, List } from "@mui/material";
import React from "react";

import RuleListItem from "./RuleListItem";

import { useScopeQuery } from "lib/graphql/generated";

function Rules(): JSX.Element {
  const { loading, error, data } = useScopeQuery();

  return (
    <div>
      {loading && <CircularProgress />}
      {error && <Alert severity="error">Error fetching scope: {error.message}</Alert>}
      {data && data.scope.length > 0 && (
        <List
          sx={{
            bgcolor: "background.paper",
          }}
        >
          {data.scope.map((rule, index) => (
            <RuleListItem key={index} rule={rule} scope={data.scope} index={index} />
          ))}
        </List>
      )}
    </div>
  );
}

export default Rules;
