import { gql, useMutation, useQuery } from "@apollo/client";
import {
  Avatar,
  Box,
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  IconButton,
  List,
  ListItem,
  ListItemAvatar,
  ListItemSecondaryAction,
  ListItemText,
  Paper,
  Snackbar,
  Tooltip,
  Typography,
  useTheme,
} from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import DescriptionIcon from "@mui/icons-material/Description";
import DeleteIcon from "@mui/icons-material/Delete";
import LaunchIcon from "@mui/icons-material/Launch";
import { Alert } from "@mui/lab";
import React, { useState } from "react";

import { Project } from "../../lib/Project";

const PROJECTS = gql`
  query Projects {
    projects {
      id
      name
      isActive
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

const CLOSE_PROJECT = gql`
  mutation CloseProject {
    closeProject {
      success
    }
  }
`;

const DELETE_PROJECT = gql`
  mutation DeleteProject($id: ID!) {
    deleteProject(id: $id) {
      success
    }
  }
`;

function ProjectList(): JSX.Element {
  const theme = useTheme();
  const {
    loading: projLoading,
    error: projErr,
    data: projData,
  } = useQuery<{ projects: Project[] }>(PROJECTS, {
    fetchPolicy: "network-only",
  });
  const [openProject, { error: openProjErr, loading: openProjLoading }] = useMutation<{ openProject: Project }>(
    OPEN_PROJECT,
    {
      errorPolicy: "all",
      onError: () => {},
      update(cache, { data }) {
        cache.modify({
          fields: {
            activeProject() {
              const activeProjRef = cache.writeFragment({
                data: data?.openProject,
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
                id: data?.openProject.id,
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
            httpRequestLogFilter(_, { DELETE }) {
              return DELETE;
            },
          },
        });
      },
    }
  );
  const [closeProject, { error: closeProjErr, client }] = useMutation(CLOSE_PROJECT, {
    errorPolicy: "all",
    onError: () => {},
    onCompleted() {
      client.resetStore();
    },
    update(cache) {
      cache.modify({
        fields: {
          activeProject() {
            return null;
          },
          projects(_, { DELETE }) {
            return DELETE;
          },
          httpRequestLogFilter(_, { DELETE }) {
            return DELETE;
          },
        },
      });
    },
  });
  const [deleteProject, { loading: deleteProjLoading, error: deleteProjErr }] = useMutation(DELETE_PROJECT, {
    errorPolicy: "all",
    onError: () => {},
    update(cache) {
      cache.modify({
        fields: {
          projects(_, { DELETE }) {
            return DELETE;
          },
        },
      });
      setDeleteDiagOpen(false);
      setDeleteNotifOpen(true);
    },
  });

  const [deleteProj, setDeleteProj] = useState<Project>();
  const [deleteDiagOpen, setDeleteDiagOpen] = useState(false);
  const handleDeleteButtonClick = (project: any) => {
    setDeleteProj(project);
    setDeleteDiagOpen(true);
  };
  const handleDeleteConfirm = () => {
    deleteProject({ variables: { id: deleteProj?.id } });
  };
  const handleDeleteCancel = () => {
    setDeleteDiagOpen(false);
  };

  const [deleteNotifOpen, setDeleteNotifOpen] = useState(false);
  const handleCloseDeleteNotif = (_: Event | React.SyntheticEvent, reason?: string) => {
    if (reason === "clickaway") {
      return;
    }
    setDeleteNotifOpen(false);
  };

  return (
    <div>
      <Dialog open={deleteDiagOpen} onClose={handleDeleteCancel}>
        <DialogTitle>
          Delete project “<strong>{deleteProj?.name}</strong>”?
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Deleting a project permanently removes all its data from the database. This action is irreversible.
          </DialogContentText>
          {deleteProjErr && <Alert severity="error">Error closing project: {deleteProjErr.message}</Alert>}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDeleteCancel} autoFocus color="secondary" variant="contained">
            Cancel
          </Button>
          <Button
            sx={{
              color: "white",
              backgroundColor: "error.main",
              "&:hover": {
                backgroundColor: "error.dark",
              },
            }}
            onClick={handleDeleteConfirm}
            disabled={deleteProjLoading}
            variant="contained"
          >
            Delete
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={deleteNotifOpen}
        autoHideDuration={3000}
        onClose={handleCloseDeleteNotif}
        anchorOrigin={{ horizontal: "center", vertical: "bottom" }}
      >
        <Alert onClose={handleCloseDeleteNotif} severity="info">
          Project <strong>{deleteProj?.name}</strong> was deleted.
        </Alert>
      </Snackbar>

      <Box mb={3}>
        <Typography variant="h6">Manage projects</Typography>
      </Box>

      <Box mb={4}>
        {projLoading && <CircularProgress />}
        {projErr && <Alert severity="error">Error fetching projects: {projErr.message}</Alert>}
        {openProjErr && <Alert severity="error">Error opening project: {openProjErr.message}</Alert>}
        {closeProjErr && <Alert severity="error">Error closing project: {closeProjErr.message}</Alert>}
      </Box>

      {projData && projData.projects.length > 0 && (
        <Paper>
          <List>
            {projData.projects.map((project) => (
              <ListItem key={project.id}>
                <ListItemAvatar>
                  <Avatar
                    sx={{
                      ...(project.isActive && {
                        color: theme.palette.secondary.dark,
                        backgroundColor: theme.palette.primary.main,
                      }),
                    }}
                  >
                    <DescriptionIcon />
                  </Avatar>
                </ListItemAvatar>
                <ListItemText>
                  {project.name} {project.isActive && <em>(Active)</em>}
                </ListItemText>
                <ListItemSecondaryAction>
                  {project.isActive && (
                    <Tooltip title="Close project">
                      <IconButton onClick={() => closeProject()}>
                        <CloseIcon />
                      </IconButton>
                    </Tooltip>
                  )}
                  {!project.isActive && (
                    <Tooltip title="Open project">
                      <span>
                        <IconButton
                          disabled={openProjLoading || projLoading}
                          onClick={() =>
                            openProject({
                              variables: { id: project.id },
                            })
                          }
                        >
                          <LaunchIcon />
                        </IconButton>
                      </span>
                    </Tooltip>
                  )}
                  <Tooltip title="Delete project">
                    <span>
                      <IconButton onClick={() => handleDeleteButtonClick(project)} disabled={project.isActive}>
                        <DeleteIcon />
                      </IconButton>
                    </span>
                  </Tooltip>
                </ListItemSecondaryAction>
              </ListItem>
            ))}
          </List>
        </Paper>
      )}
      {projData?.projects.length === 0 && (
        <Alert severity="info">There are no projects. Create one to get started.</Alert>
      )}
    </div>
  );
}

export default ProjectList;
