import { gql, useQuery } from "@apollo/client";
import { CircularProgress, List } from "@mui/material";
import { Alert } from "@mui/lab";
import React from "react";
import RuleListItem from "./RuleListItem";
import { ScopeRule } from "../../lib/scope";

export const SCOPE = gql`
  query Scope {
    scope {
      url
    }
  }
`;

function Rules(): JSX.Element {
  const { loading, error, data } = useQuery<{ scope: ScopeRule[] }>(SCOPE);

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
