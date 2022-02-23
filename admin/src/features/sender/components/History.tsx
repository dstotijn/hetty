import { TableContainer, Table, TableHead, TableRow, TableCell, Typography, Box, TableBody } from "@mui/material";
import { useRouter } from "next/router";

import CenteredPaper from "lib/components/CenteredPaper";
import HttpStatusIcon from "lib/components/HttpStatusIcon";
import { useGetSenderRequestsQuery } from "lib/graphql/generated";

function History(): JSX.Element {
  const { data, loading } = useGetSenderRequestsQuery({
    pollInterval: 1000,
  });

  const router = useRouter();
  const activeId = router.query.id as string | undefined;

  const handleRowClick = (id: string) => {
    router.push(`/sender?id=${id}`);
  };

  return (
    <Box>
      {!loading && data?.senderRequests && data?.senderRequests.length > 0 && (
        <TableContainer sx={{ overflowX: "initial" }}>
          <Table size="small" stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell>Method</TableCell>
                <TableCell>Origin</TableCell>
                <TableCell>Path</TableCell>
                <TableCell>Status</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {data?.senderRequests &&
                data.senderRequests.map(({ id, method, url, response }) => {
                  const { origin, pathname, search, hash } = new URL(url);

                  const cellStyle = {
                    whiteSpace: "nowrap",
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                  };

                  return (
                    <TableRow
                      key={id}
                      sx={{
                        "&:hover": {
                          cursor: "pointer",
                        },
                        ...(id === activeId && {
                          bgcolor: "action.selected",
                          cursor: "inherit",
                        }),
                      }}
                      hover
                      onClick={() => handleRowClick(id)}
                    >
                      <TableCell sx={{ ...cellStyle, width: "100px" }}>
                        <code>{method}</code>
                      </TableCell>
                      <TableCell sx={{ ...cellStyle, maxWidth: "100px" }}>{origin}</TableCell>
                      <TableCell sx={{ ...cellStyle, maxWidth: "200px" }}>
                        {decodeURIComponent(pathname + search + hash)}
                      </TableCell>
                      <TableCell style={{ maxWidth: "100px" }}>
                        {response && (
                          <div>
                            <HttpStatusIcon status={response.statusCode} />{" "}
                            <code>
                              {response.statusCode} {response.statusReason}
                            </code>
                          </div>
                        )}
                      </TableCell>
                    </TableRow>
                  );
                })}
            </TableBody>
          </Table>
        </TableContainer>
      )}
      <Box sx={{ mt: 2, height: "100%" }}>
        {!loading && data?.senderRequests.length === 0 && (
          <CenteredPaper>
            <Typography>No requests created yet.</Typography>
          </CenteredPaper>
        )}
      </Box>
    </Box>
  );
}

export default History;
