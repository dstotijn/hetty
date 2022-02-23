import AddIcon from "@mui/icons-material/Add";
import { Box, Button, CircularProgress, TextField, Typography } from "@mui/material";
import React, { useState } from "react";

import useOpenProjectMutation from "../hooks/useOpenProjectMutation";

import { useCreateProjectMutation } from "lib/graphql/generated";

function NewProject(): JSX.Element {
  const [name, setName] = useState("");

  const [createProject, createProjResult] = useCreateProjectMutation({
    onCompleted(data) {
      setName("");
      if (data?.createProject) {
        openProject({ variables: { id: data.createProject?.id } });
      }
    },
  });
  const [openProject, openProjResult] = useOpenProjectMutation();

  const handleCreateAndOpenProjectForm = (e: React.SyntheticEvent) => {
    e.preventDefault();
    createProject({ variables: { name } });
  };

  return (
    <div>
      <Box mb={3}>
        <Typography variant="h6">New project</Typography>
      </Box>
      <form onSubmit={handleCreateAndOpenProjectForm} autoComplete="off">
        <TextField
          sx={{
            mr: 2,
          }}
          color="primary"
          size="small"
          label="Project name"
          placeholder="Project name…"
          onChange={(e) => setName(e.target.value)}
          error={Boolean(createProjResult.error || openProjResult.error)}
          helperText={
            (createProjResult.error && createProjResult.error.message) ||
            (openProjResult.error && openProjResult.error.message)
          }
        />
        <Button
          type="submit"
          variant="contained"
          color="primary"
          size="large"
          sx={{
            pt: 0.9,
            pb: 0.7,
          }}
          disabled={createProjResult.loading || openProjResult.loading}
          startIcon={createProjResult.loading || openProjResult.loading ? <CircularProgress size={22} /> : <AddIcon />}
        >
          Create & open project
        </Button>
      </form>
    </div>
  );
}

export default NewProject;
