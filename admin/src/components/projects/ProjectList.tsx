import { gql, useMutation, useQuery } from "@apollo/client";
import {
  Avatar,
  Box,
  Button,
  CircularProgress,
  createStyles,
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
  makeStyles,
  Snackbar,
  Theme,
  Tooltip,
  Typography,
} from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";
import DescriptionIcon from "@material-ui/icons/Description";
import DeleteIcon from "@material-ui/icons/Delete";
import LaunchIcon from "@material-ui/icons/Launch";
import { Alert } from "@material-ui/lab";

import React, { useState } from "react";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    projectsList: {
      backgroundColor: theme.palette.background.paper,
    },
    activeProject: {
      color: theme.palette.getContrastText(theme.palette.secondary.main),
      backgroundColor: theme.palette.secondary.main,
    },
    deleteProjectButton: {
      color: theme.palette.error.main,
    },
  })
);

const PROJECTS = gql`
  query Projects {
    projects {
      name
      isActive
    }
  }
`;

const OPEN_PROJECT = gql`
  mutation OpenProject($name: String!) {
    openProject(name: $name) {
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
  mutation DeleteProject($name: String!) {
    deleteProject(name: $name) {
      success
    }
  }
`;

function ProjectList(): JSX.Element {
  const classes = useStyles();
  const { loading: projLoading, error: projErr, data: projData } = useQuery(
    PROJECTS
  );
  const [
    openProject,
    { error: openProjErr, loading: openProjLoading },
  ] = useMutation(OPEN_PROJECT, {
    errorPolicy: "all",
    onError: () => {},
    update(cache, { data: { openProject } }) {
      cache.modify({
        fields: {
          activeProject() {
            const activeProjRef = cache.writeFragment({
              data: openProject,
              fragment: gql`
                fragment ActiveProject on Project {
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
              id: openProject.name,
              data: openProject,
              fragment: gql`
                fragment OpenProject on Project {
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
  });
  const [closeProject, { error: closeProjErr }] = useMutation(CLOSE_PROJECT, {
    errorPolicy: "all",
    onError: () => {},
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
  const [
    deleteProject,
    { loading: deleteProjLoading, error: deleteProjErr },
  ] = useMutation(DELETE_PROJECT, {
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

  const [deleteProjName, setDeleteProjName] = useState(null);
  const [deleteDiagOpen, setDeleteDiagOpen] = useState(false);
  const handleDeleteButtonClick = (name: string) => {
    setDeleteProjName(name);
    setDeleteDiagOpen(true);
  };
  const handleDeleteConfirm = () => {
    deleteProject({ variables: { name: deleteProjName } });
  };
  const handleDeleteCancel = () => {
    setDeleteDiagOpen(false);
  };

  const [deleteNotifOpen, setDeleteNotifOpen] = useState(false);
  const handleCloseDeleteNotif = (_, reason?: string) => {
    if (reason === "clickaway") {
      return;
    }
    setDeleteNotifOpen(false);
  };

  return (
    <div>
      <Dialog open={deleteDiagOpen} onClose={handleDeleteCancel}>
        <DialogTitle>
          Delete project “<strong>{deleteProjName}</strong>”?
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Deleting a project permanently removes its database file from disk.
            This action is irreversible.
          </DialogContentText>
          {deleteProjErr && (
            <Alert severity="error">
              Error closing project: {deleteProjErr.message}
            </Alert>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDeleteCancel} autoFocus>
            Cancel
          </Button>
          <Button
            className={classes.deleteProjectButton}
            onClick={handleDeleteConfirm}
            disabled={deleteProjLoading}
          >
            Delete
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={deleteNotifOpen}
        autoHideDuration={3000}
        onClose={handleCloseDeleteNotif}
      >
        <Alert onClose={handleCloseDeleteNotif} severity="info">
          Project <strong>{deleteProjName}</strong> was deleted.
        </Alert>
      </Snackbar>

      <Box mb={3}>
        <Typography variant="h6">Manage projects</Typography>
      </Box>

      <Box mb={4}>
        {projLoading && <CircularProgress />}
        {projErr && (
          <Alert severity="error">
            Error fetching projects: {projErr.message}
          </Alert>
        )}
        {openProjErr && (
          <Alert severity="error">
            Error opening project: {openProjErr.message}
          </Alert>
        )}
        {closeProjErr && (
          <Alert severity="error">
            Error closing project: {closeProjErr.message}
          </Alert>
        )}
      </Box>

      {projData?.projects.length > 0 && (
        <List className={classes.projectsList}>
          {projData.projects.map((project) => (
            <ListItem key={project.name}>
              <ListItemAvatar>
                <Avatar
                  className={
                    project.isActive ? classes.activeProject : undefined
                  }
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
                            variables: { name: project.name },
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
                    <IconButton
                      onClick={() => handleDeleteButtonClick(project.name)}
                      disabled={project.isActive}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </span>
                </Tooltip>
              </ListItemSecondaryAction>
            </ListItem>
          ))}
        </List>
      )}
      {projData?.projects.length === 0 && (
        <Alert severity="info">
          There are no projects. Create one to get started.
        </Alert>
      )}
    </div>
  );
}

export default ProjectList;
