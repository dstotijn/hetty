import DeleteIcon from "@mui/icons-material/Delete";
import FilterListIcon from "@mui/icons-material/FilterList";
import SearchIcon from "@mui/icons-material/Search";
import { Alert } from "@mui/lab";
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
import React, { useRef, useState } from "react";

import { ConfirmationDialog, useConfirmationDialog } from "lib/components/ConfirmationDialog";
import {
  HttpRequestLogFilterDocument,
  HttpRequestLogsDocument,
  useClearHttpRequestLogMutation,
  useHttpRequestLogFilterQuery,
  useSetHttpRequestLogFilterMutation,
} from "lib/graphql/generated";
import { withoutTypename } from "lib/graphql/omitTypename";

function Search(): JSX.Element {
  const theme = useTheme();

  const [searchExpr, setSearchExpr] = useState("");
  const filterResult = useHttpRequestLogFilterQuery({
    onCompleted: (data) => {
      setSearchExpr(data.httpRequestLogFilter?.searchExpression || "");
    },
  });
  const filter = filterResult.data?.httpRequestLogFilter;

  const [setFilterMutate, setFilterResult] = useSetHttpRequestLogFilterMutation({
    update(cache, { data }) {
      cache.writeQuery({
        query: HttpRequestLogFilterDocument,
        data: {
          httpRequestLogFilter: data?.setHttpRequestLogFilter,
        },
      });
    },
  });

  const [clearHTTPRequestLog, clearHTTPRequestLogResult] = useClearHttpRequestLogMutation({
    refetchQueries: [{ query: HttpRequestLogsDocument }],
  });
  const clearHTTPConfirmationDialog = useConfirmationDialog();

  const filterRef = useRef<HTMLFormElement>(null);
  const [filterOpen, setFilterOpen] = useState(false);

  const handleSubmit = (e: React.SyntheticEvent) => {
    setFilterMutate({
      variables: {
        filter: {
          ...withoutTypename(filter),
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
      <Error prefix="Error fetching filter" error={filterResult.error} />
      <Error prefix="Error setting filter" error={setFilterResult.error} />
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
                  color: filter?.onlyInScope ? "primary.main" : "inherit",
                }}
              >
                {filterResult.loading || setFilterResult.loading ? (
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
                      checked={filter?.onlyInScope ? true : false}
                      disabled={filterResult.loading || setFilterResult.loading}
                      onChange={(e) =>
                        setFilterMutate({
                          variables: {
                            filter: {
                              ...withoutTypename(filter),
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
