import { gql, useMutation } from "@apollo/client";
import { Box, Button, CircularProgress, TextField, Typography } from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import React, { useState } from "react";

const CREATE_PROJECT = gql`
  mutation CreateProject($name: String!) {
    createProject(name: $name) {
      id
      name
    }
  }
`;

const OPEN_PROJECT = gql`
  mutation OpenProject($id: ID!) {
    openProject(id: $id) {
      id
      name
      isActive
    }
  }
`;

function NewProject(): JSX.Element {
  const [name, setName] = useState("");

  const [createProject, { error: createProjErr, loading: createProjLoading }] = useMutation(CREATE_PROJECT, {
    onError: () => {},
    onCompleted(data) {
      setName("");
      openProject({ variables: { id: data.createProject.id } });
    },
  });
  const [openProject, { error: openProjErr, loading: openProjLoading }] = useMutation(OPEN_PROJECT, {
    onError: () => {},
    update(cache, { data: { openProject } }) {
      cache.modify({
        fields: {
          activeProject() {
            const activeProjRef = cache.writeFragment({
              id: openProject.id,
              data: openProject,
              fragment: gql`
                fragment ActiveProject on Project {
                  id
                  name
                  isActive
                  type
                }
              `,
            });
            return activeProjRef;
          },
          projects(_, { DELETE }) {
            cache.writeFragment({
              id: openProject.id,
              data: openProject,
              fragment: gql`
                fragment OpenProject on Project {
                  id
                  name
                  isActive
                  type
                }
              `,
            });
            return DELETE;
          },
        },
      });
    },
  });

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
          error={Boolean(createProjErr || openProjErr)}
          helperText={(createProjErr && createProjErr.message) || (openProjErr && openProjErr.message)}
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
          disabled={createProjLoading || openProjLoading}
          startIcon={createProjLoading || openProjLoading ? <CircularProgress size={22} /> : <AddIcon />}
        >
          Create & open project
        </Button>
      </form>
    </div>
  );
}

export default NewProject;
