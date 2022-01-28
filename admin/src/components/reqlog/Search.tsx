import {
  Box,
  Checkbox,
  CircularProgress,
  ClickAwayListener,
  FormControlLabel,
  InputBase,
  Paper,
  Popper,
  Tooltip,
  useTheme,
} from "@mui/material";
import IconButton from "@mui/material/IconButton";
import SearchIcon from "@mui/icons-material/Search";
import FilterListIcon from "@mui/icons-material/FilterList";
import DeleteIcon from "@mui/icons-material/Delete";
import React, { useRef, useState } from "react";
import { gql, useMutation, useQuery } from "@apollo/client";
import { withoutTypename } from "../../lib/omitTypename";
import { Alert } from "@mui/lab";
import { useClearHTTPRequestLog } from "./hooks/useClearHTTPRequestLog";
import { ConfirmationDialog, useConfirmationDialog } from "./ConfirmationDialog";

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

export interface SearchFilter {
  onlyInScope: boolean;
  searchExpression: string;
}

function Search(): JSX.Element {
  const theme = useTheme();

  const [searchExpr, setSearchExpr] = useState("");
  const {
    loading: filterLoading,
    error: filterErr,
    data: filter,
  } = useQuery(FILTER, {
    onCompleted: (data) => {
      setSearchExpr(data.httpRequestLogFilter?.searchExpression || "");
    },
  });

  const [setFilterMutate, { error: setFilterErr, loading: setFilterLoading }] = useMutation<{
    setHttpRequestLogFilter: SearchFilter | null;
  }>(SET_FILTER, {
    update(cache, { data }) {
      cache.writeQuery({
        query: FILTER,
        data: {
          httpRequestLogFilter: data?.setHttpRequestLogFilter,
        },
      });
    },
    onError: () => {},
  });

  const [clearHTTPRequestLog, clearHTTPRequestLogResult] = useClearHTTPRequestLog();
  const clearHTTPConfirmationDialog = useConfirmationDialog();

  const filterRef = useRef<HTMLFormElement>(null);
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

  const handleClickAway = (event: MouseEvent | TouchEvent) => {
    if (filterRef?.current && filterRef.current.contains(event.target as HTMLElement)) {
      return;
    }
    setFilterOpen(false);
  };

  return (
    <Box>
      <Error prefix="Error fetching filter" error={filterErr} />
      <Error prefix="Error setting filter" error={setFilterErr} />
      <Error prefix="Error clearing all HTTP logs" error={clearHTTPRequestLogResult.error} />
      <Box style={{ display: "flex", flex: 1 }}>
        <ClickAwayListener onClickAway={handleClickAway}>
          <Paper
            component="form"
            onSubmit={handleSubmit}
            ref={filterRef}
            sx={{
              padding: "2px 4px",
              display: "flex",
              alignItems: "center",
              width: 400,
            }}
          >
            <Tooltip title="Toggle filter options">
              <IconButton
                onClick={() => setFilterOpen(!filterOpen)}
                sx={{
                  p: 1,
                  color: filter?.httpRequestLogFilter?.onlyInScope ? "primary.main" : "inherit",
                }}
              >
                {filterLoading || setFilterLoading ? (
                  <CircularProgress sx={{ color: theme.palette.text.primary }} size={23} />
                ) : (
                  <FilterListIcon />
                )}
              </IconButton>
            </Tooltip>
            <InputBase
              sx={{
                ml: 1,
                flex: 1,
              }}
              placeholder="Search proxy logs…"
              value={searchExpr}
              onChange={(e) => setSearchExpr(e.target.value)}
              onFocus={() => setFilterOpen(true)}
            />
            <Tooltip title="Search">
              <IconButton type="submit" sx={{ padding: 1.25 }}>
                <SearchIcon />
              </IconButton>
            </Tooltip>
            <Popper
              open={filterOpen}
              anchorEl={filterRef.current}
              placement="bottom"
              style={{ zIndex: theme.zIndex.appBar }}
            >
              <Paper
                sx={{
                  width: 400,
                  marginTop: 0.5,
                  p: 1.5,
                }}
              >
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={filter?.httpRequestLogFilter?.onlyInScope ? true : false}
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
