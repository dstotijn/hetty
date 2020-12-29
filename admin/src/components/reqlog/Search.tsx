import {
  Box,
  Checkbox,
  CircularProgress,
  ClickAwayListener,
  createStyles,
  FormControlLabel,
  InputBase,
  makeStyles,
  Paper,
  Popper,
  Theme,
  Tooltip,
  useTheme,
} from "@material-ui/core";
import IconButton from "@material-ui/core/IconButton";
import SearchIcon from "@material-ui/icons/Search";
import FilterListIcon from "@material-ui/icons/FilterList";
import DeleteIcon from "@material-ui/icons/Delete";
import React, { useRef, useState } from "react";
import { gql, useMutation, useQuery } from "@apollo/client";
import { withoutTypename } from "../../lib/omitTypename";
import { Alert } from "@material-ui/lab";
import { useClearHTTPRequestLog } from "./hooks/useClearHTTPRequestLog";
import {
  ConfirmationDialog,
  useConfirmationDialog,
} from "./ConfirmationDialog";

const FILTER = gql`
  query HttpRequestLogFilter {
    httpRequestLogFilter {
      onlyInScope
      searchExpression
    }
  }
`;

const SET_FILTER = gql`
  mutation SetHttpRequestLogFilter($filter: HttpRequestLogFilterInput) {
    setHttpRequestLogFilter(filter: $filter) {
      onlyInScope
      searchExpression
    }
  }
`;

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      padding: "2px 4px",
      display: "flex",
      alignItems: "center",
      width: 400,
    },
    input: {
      marginLeft: theme.spacing(1),
      flex: 1,
    },
    iconButton: {
      padding: 10,
    },
    filterPopper: {
      width: 400,
      marginTop: 6,
      zIndex: 99,
    },
    filterOptions: {
      padding: theme.spacing(2),
    },
    filterLoading: {
      marginRight: 1,
      color: theme.palette.text.primary,
    },
  })
);

export interface SearchFilter {
  onlyInScope: boolean;
  searchExpression: string;
}

function Search(): JSX.Element {
  const classes = useStyles();
  const theme = useTheme();

  const [searchExpr, setSearchExpr] = useState("");
  const { loading: filterLoading, error: filterErr, data: filter } = useQuery(
    FILTER,
    {
      onCompleted: (data) => {
        setSearchExpr(data.httpRequestLogFilter?.searchExpression || "");
      },
    }
  );

  const [
    setFilterMutate,
    { error: setFilterErr, loading: setFilterLoading },
  ] = useMutation<{
    setHttpRequestLogFilter: SearchFilter | null;
  }>(SET_FILTER, {
    update(cache, { data: { setHttpRequestLogFilter } }) {
      cache.writeQuery({
        query: FILTER,
        data: {
          httpRequestLogFilter: setHttpRequestLogFilter,
        },
      });
    },
    onError: () => {},
  });

  const [
    clearHTTPRequestLog,
    clearHTTPRequestLogResult,
  ] = useClearHTTPRequestLog();
  const clearHTTPConfirmationDialog = useConfirmationDialog();

  const filterRef = useRef<HTMLElement | null>();
  const [filterOpen, setFilterOpen] = useState(false);

  const handleSubmit = (e: React.SyntheticEvent) => {
    setFilterMutate({
      variables: {
        filter: {
          ...withoutTypename(filter?.httpRequestLogFilter),
          searchExpression: searchExpr,
        },
      },
    });
    setFilterOpen(false);
    e.preventDefault();
  };

  const handleClickAway = (event: React.MouseEvent<EventTarget>) => {
    if (filterRef.current.contains(event.target as HTMLElement)) {
      return;
    }
    setFilterOpen(false);
  };

  return (
    <Box>
      <Error prefix="Error fetching filter" error={filterErr} />
      <Error prefix="Error setting filter" error={setFilterErr} />
      <Error
        prefix="Error clearing all HTTP logs"
        error={clearHTTPRequestLogResult.error}
      />
      <Box style={{ display: "flex", flex: 1 }}>
        <ClickAwayListener onClickAway={handleClickAway}>
          <Paper
            component="form"
            onSubmit={handleSubmit}
            ref={filterRef}
            className={classes.root}
          >
            <Tooltip title="Toggle filter options">
              <IconButton
                className={classes.iconButton}
                onClick={() => setFilterOpen(!filterOpen)}
                style={{
                  color: filter?.httpRequestLogFilter?.onlyInScope
                    ? theme.palette.secondary.main
                    : "inherit",
                }}
              >
                {filterLoading || setFilterLoading ? (
                  <CircularProgress
                    className={classes.filterLoading}
                    size={23}
                  />
                ) : (
                  <FilterListIcon />
                )}
              </IconButton>
            </Tooltip>
            <InputBase
              className={classes.input}
              placeholder="Search proxy logs…"
              value={searchExpr}
              onChange={(e) => setSearchExpr(e.target.value)}
              onFocus={() => setFilterOpen(true)}
            />
            <Tooltip title="Search">
              <IconButton type="submit" className={classes.iconButton}>
                <SearchIcon />
              </IconButton>
            </Tooltip>
            <Popper
              className={classes.filterPopper}
              open={filterOpen}
              anchorEl={filterRef.current}
              placement="bottom-start"
            >
              <Paper className={classes.filterOptions}>
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={
                        filter?.httpRequestLogFilter?.onlyInScope ? true : false
                      }
                      disabled={filterLoading || setFilterLoading}
                      onChange={(e) =>
                        setFilterMutate({
                          variables: {
                            filter: {
                              ...withoutTypename(filter?.httpRequestLogFilter),
                              onlyInScope: e.target.checked,
                            },
                          },
                        })
                      }
                    />
                  }
                  label="Only show in-scope requests"
                />
              </Paper>
            </Popper>
          </Paper>
        </ClickAwayListener>
        <Box style={{ marginLeft: "auto" }}>
          <Tooltip title="Clear all">
            <IconButton onClick={clearHTTPConfirmationDialog.open}>
              <DeleteIcon />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>
      <ConfirmationDialog
        isOpen={clearHTTPConfirmationDialog.isOpen}
        onClose={clearHTTPConfirmationDialog.close}
        onConfirm={clearHTTPRequestLog}
      >
        All proxy logs are going to be removed. This action cannot be undone.
      </ConfirmationDialog>
    </Box>
  );
}

function Error(props: { prefix: string; error?: Error }) {
  if (!props.error) return null;

  return (
    <Box mb={4}>
      <Alert severity="error">
        {props.prefix}: {props.error.message}
      </Alert>
    </Box>
  );
}

export default Search;
