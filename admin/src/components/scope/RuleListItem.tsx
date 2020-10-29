import { gql, useApolloClient, useMutation, useQuery } from "@apollo/client";
import {
  Avatar,
  Chip,
  IconButton,
  ListItem,
  ListItemAvatar,
  ListItemSecondaryAction,
  ListItemText,
  Tooltip,
} from "@material-ui/core";
import CodeIcon from "@material-ui/icons/Code";
import DeleteIcon from "@material-ui/icons/Delete";
import React from "react";
import { SCOPE } from "./Rules";

const SET_SCOPE = gql`
  mutation SetScope($scope: [ScopeRuleInput!]!) {
    setScope(scope: $scope) {
      url
    }
  }
`;

function RuleListItem({ scope, rule, index }): JSX.Element {
  const client = useApolloClient();
  const [setScope, { loading }] = useMutation(SET_SCOPE, {
    update(_, { data: { setScope } }) {
      client.writeQuery({
        query: SCOPE,
        data: { scope: setScope },
      });
    },
  });

  const handleDelete = (index: number) => {
    const clone = [...scope];
    clone.splice(index, 1);
    setScope({
      variables: {
        scope: clone.map(({ url }) => ({ url })),
      },
    });
  };

  return (
    <ListItem>
      <ListItemAvatar>
        <Avatar>
          <CodeIcon />
        </Avatar>
      </ListItemAvatar>
      <RuleListItemText rule={rule} />
      <ListItemSecondaryAction>
        <RuleTypeChip rule={rule} />
        <Tooltip title="Delete rule">
          <span style={{ marginLeft: 8 }}>
            <IconButton onClick={() => handleDelete(index)} disabled={loading}>
              <DeleteIcon />
            </IconButton>
          </span>
        </Tooltip>
      </ListItemSecondaryAction>
    </ListItem>
  );
}

function RuleListItemText({ rule }): JSX.Element {
  let text: JSX.Element;

  if (rule.url) {
    text = <code>{rule.url}</code>;
  }

  // TODO: Parse and handle rule.header and rule.body.

  return <ListItemText>{text}</ListItemText>;
}

function RuleTypeChip({ rule }): JSX.Element {
  if (rule.url) {
    return <Chip label="URL" variant="outlined" />;
  }
}

export default RuleListItem;
