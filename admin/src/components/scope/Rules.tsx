import { gql, useQuery } from "@apollo/client";
import {
  CircularProgress,
  createStyles,
  List,
  makeStyles,
  Theme,
} from "@material-ui/core";
import { Alert } from "@material-ui/lab";
import React from "react";
import RuleListItem from "./RuleListItem";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    rulesList: {
      backgroundColor: theme.palette.background.paper,
    },
  })
);

export const SCOPE = gql`
  query Scope {
    scope {
      url
    }
  }
`;

function Rules(): JSX.Element {
  const classes = useStyles();
  const { loading, error, data } = useQuery(SCOPE);

  return (
    <div>
      {loading && <CircularProgress />}
      {error && (
        <Alert severity="error">Error fetching scope: {error.message}</Alert>
      )}
      {data?.scope.length > 0 && (
        <List className={classes.rulesList}>
          {data.scope.map((rule, index) => (
            <RuleListItem
              key={index}
              rule={rule}
              scope={data.scope}
              index={index}
            />
          ))}
        </List>
      )}
    </div>
  );
}

export default Rules;
