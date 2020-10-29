import { gql, useApolloClient, useMutation } from "@apollo/client";
import {
  Box,
  Button,
  CircularProgress,
  createStyles,
  FormControl,
  FormControlLabel,
  FormLabel,
  makeStyles,
  Radio,
  RadioGroup,
  TextField,
  Theme,
} from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import { Alert } from "@material-ui/lab";
import React from "react";
import { SCOPE } from "./Rules";

const SET_SCOPE = gql`
  mutation SetScope($scope: [ScopeRuleInput!]!) {
    setScope(scope: $scope) {
      url
    }
  }
`;

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    ruleExpression: {
      fontFamily: "'JetBrains Mono', monospace",
    },
  })
);

function AddRule(): JSX.Element {
  const classes = useStyles();

  const [ruleType, setRuleType] = React.useState("url");
  const [expression, setExpression] = React.useState(null);

  const client = useApolloClient();
  const [setScope, { error, loading }] = useMutation(SET_SCOPE, {
    onError() {},
    onCompleted() {
      expression.value = "";
    },
    update(_, { data: { setScope } }) {
      client.writeQuery({
        query: SCOPE,
        data: { scope: setScope },
      });
    },
  });

  const handleTypeChange = (e: React.ChangeEvent, value: string) => {
    setRuleType(value);
  };
  const handleSubmit = (e: React.SyntheticEvent) => {
    e.preventDefault();
    let scope = [];

    try {
      const data = client.readQuery({
        query: SCOPE,
      });
      scope = data.scope;
    } catch (e) {}

    setScope({
      variables: {
        scope: [
          ...scope.map(({ url }) => ({ url })),
          { url: expression.value },
        ],
      },
    });
  };

  return (
    <div>
      {error && (
        <Box mb={4}>
          <Alert severity="error">Error adding rule: {error.message}</Alert>
        </Box>
      )}
      <form onSubmit={handleSubmit} autoComplete="off">
        <FormControl fullWidth>
          <FormLabel color="secondary" component="legend">
            Rule Type
          </FormLabel>
          <RadioGroup
            row
            name="ruleType"
            value={ruleType}
            onChange={handleTypeChange}
          >
            <FormControlLabel value="url" control={<Radio />} label="URL" />
          </RadioGroup>
        </FormControl>
        <FormControl fullWidth>
          <TextField
            label="Expression"
            placeholder="^https:\/\/(.*)example.com(.*)"
            helperText="Regular expression to match on."
            color="secondary"
            variant="outlined"
            required
            InputProps={{
              className: classes.ruleExpression,
            }}
            InputLabelProps={{
              shrink: true,
            }}
            inputProps={{
              ref: (node) => {
                setExpression(node);
              },
            }}
            margin="normal"
          />
        </FormControl>
        <Box my={2}>
          <Button
            type="submit"
            variant="contained"
            color="secondary"
            disabled={loading}
            startIcon={loading ? <CircularProgress size={22} /> : <AddIcon />}
          >
            Add rule
          </Button>
        </Box>
      </form>
    </div>
  );
}

export default AddRule;
