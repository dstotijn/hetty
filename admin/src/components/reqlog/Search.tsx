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
import React, { useRef, useState } from "react";
import { gql, useApolloClient, useMutation, useQuery } from "@apollo/client";
import { withoutTypename } from "../../lib/omitTypename";
import { Alert } from "@material-ui/lab";

const FILTER = gql`
  query HttpRequestLogFilter {
    httpRequestLogFilter {
      onlyInScope
    }
  }
`;

const SET_FILTER = gql`
  mutation SetHttpRequestLogFilter($filter: HttpRequestLogFilterInput) {
    setHttpRequestLogFilter(filter: $filter) {
      onlyInScope
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
}

function Search(): JSX.Element {
  const classes = useStyles();
  const theme = useTheme();

  const { loading: filterLoading, error: filterErr, data: filter } = useQuery(
    FILTER
  );

  const client = useApolloClient();
  const [
    setFilterMutate,
    { error: setFilterErr, loading: setFilterLoading },
  ] = useMutation<{
    setHttpRequestLogFilter: SearchFilter | null;
  }>(SET_FILTER, {
    update(_, { data: { setHttpRequestLogFilter } }) {
      client.writeQuery({
        query: FILTER,
        data: {
          httpRequestLogFilter: setHttpRequestLogFilter,
        },
      });
    },
  });

  const filterRef = useRef<HTMLElement | null>();
  const [filterOpen, setFilterOpen] = useState(false);

  const handleSubmit = (e: React.SyntheticEvent) => {
    e.preventDefault();
  };

  const handleClickAway = (event: React.MouseEvent<EventTarget>) => {
    if (filterRef.current.contains(event.target as HTMLElement)) {
      return;
    }
    setFilterOpen(false);
  };

  return (
    <ClickAwayListener onClickAway={handleClickAway}>
      <Box style={{ display: "inline-block" }}>
        {filterErr && (
          <Box mb={4}>
            <Alert severity="error">
              Error fetching filter: {filterErr.message}
            </Alert>
          </Box>
        )}
        {setFilterErr && (
          <Box mb={4}>
            <Alert severity="error">
              Error setting filter: {setFilterErr.message}
            </Alert>
          </Box>
        )}
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
                color:
                  filter?.httpRequestLogFilter !== null
                    ? theme.palette.secondary.main
                    : "inherit",
              }}
            >
              {filterLoading || setFilterLoading ? (
                <CircularProgress className={classes.filterLoading} size={23} />
              ) : (
                <FilterListIcon />
              )}
            </IconButton>
          </Tooltip>
          <InputBase
            className={classes.input}
            placeholder="Search proxy logs…"
            onFocus={() => setFilterOpen(true)}
          />
          <Tooltip title="Search">
            <IconButton type="submit" className={classes.iconButton}>
              <SearchIcon />
            </IconButton>
          </Tooltip>
        </Paper>

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
      </Box>
    </ClickAwayListener>
  );
}

export default Search;
