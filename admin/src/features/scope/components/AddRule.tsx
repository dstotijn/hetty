import { useApolloClient } from "@apollo/client";
import AddIcon from "@mui/icons-material/Add";
import { Alert } from "@mui/lab";
import {
  Box,
  Button,
  CircularProgress,
  FormControl,
  FormControlLabel,
  FormLabel,
  Radio,
  RadioGroup,
  TextField,
} from "@mui/material";
import React, { useState } from "react";

import { ScopeDocument, ScopeQuery, ScopeRule, useSetScopeMutation } from "lib/graphql/generated";

function AddRule(): JSX.Element {
  const [ruleType, setRuleType] = useState("url");
  const [expression, setExpression] = useState("");

  const client = useApolloClient();
  const [setScope, { error, loading }] = useSetScopeMutation({
    onCompleted({ setScope }) {
      client.writeQuery({
        query: ScopeDocument,
        data: { scope: setScope },
      });
      setExpression("");
    },
  });

  const handleTypeChange = (e: React.ChangeEvent, value: string) => {
    setRuleType(value);
  };
  const handleSubmit = (e: React.SyntheticEvent) => {
    e.preventDefault();
    let scope: ScopeRule[] = [];

    try {
      const data = client.readQuery<ScopeQuery>({
        query: ScopeDocument,
      });
      if (data) {
        scope = data.scope;
      }
    } catch (e) {}

    setScope({
      variables: {
        scope: [...scope.map(({ url }) => ({ url })), { url: expression }],
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
          <FormLabel color="primary" component="legend">
            Rule Type
          </FormLabel>
          <RadioGroup row name="ruleType" value={ruleType} onChange={handleTypeChange}>
            <FormControlLabel value="url" control={<Radio />} label="URL" />
          </RadioGroup>
        </FormControl>
        <FormControl fullWidth>
          <TextField
            label="Expression"
            placeholder="^https:\/\/(.*)example.com(.*)"
            helperText="Regular expression to match on."
            color="primary"
            variant="outlined"
            required
            value={expression}
            onChange={(e) => setExpression(e.target.value)}
            InputProps={{
              sx: { fontFamily: "'JetBrains Mono', monospace" },
            }}
            InputLabelProps={{
              shrink: true,
            }}
            margin="normal"
          />
        </FormControl>
        <Box my={2}>
          <Button
            type="submit"
            variant="contained"
            color="primary"
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
