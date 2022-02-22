import { IconButton, InputBase, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import ClearIcon from "@mui/icons-material/Clear";

export type KeyValuePair = {
  key: string;
  value: string;
};

export type KeyValuePairTableProps = {
  items: KeyValuePair[];
  onChange?: (key: string, value: string, index: number) => void;
  onDelete?: (index: number) => void;
};

export function KeyValuePairTable({ items, onChange, onDelete }: KeyValuePairTableProps): JSX.Element {
  const inputSx = {
    fontSize: "0.875rem",
    "&.MuiInputBase-root input": {
      p: 0,
    },
  };

  return (
    <TableContainer>
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell>Key</TableCell>
            <TableCell>Value</TableCell>
            {onDelete && <TableCell padding="checkbox"></TableCell>}
          </TableRow>
        </TableHead>
        <TableBody
          sx={{
            "td, th, input": {
              fontFamily: "'JetBrains Mono', monospace",
              fontSize: "0.75rem",
              py: 0.2,
            },
            "td span, th span": {
              display: "block",
              py: 0.7,
            },
          }}
        >
          {items.map(({ key, value }, idx) => (
            <TableRow
              key={idx}
              hover
              sx={{
                "& .delete-button": {
                  visibility: "hidden",
                },
                "&:hover .delete-button": {
                  visibility: "inherit",
                },
              }}
            >
              <TableCell component="th" scope="row">
                {!onChange && <span>{key}</span>}
                {onChange && (
                  <InputBase
                    size="small"
                    fullWidth
                    placeholder="Key"
                    value={key}
                    onChange={(e) => {
                      onChange && onChange(e.target.value, value, idx);
                    }}
                    sx={inputSx}
                  />
                )}
              </TableCell>
              <TableCell sx={{ width: "60%", wordBreak: "break-all" }}>
                {!onChange && value}
                {onChange && (
                  <InputBase
                    size="small"
                    fullWidth
                    placeholder="Value"
                    value={value}
                    onChange={(e) => {
                      onChange && onChange(key, e.target.value, idx);
                    }}
                    sx={inputSx}
                  />
                )}
              </TableCell>
              {onDelete && (
                <TableCell>
                  <div className="delete-button">
                    <IconButton
                      size="small"
                      onClick={() => {
                        onDelete && onDelete(idx);
                      }}
                      sx={{
                        visibility: onDelete === undefined || items.length === idx + 1 ? "hidden" : "inherit",
                      }}
                    >
                      <ClearIcon fontSize="inherit" />
                    </IconButton>
                  </div>
                </TableCell>
              )}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

export function sortKeyValuePairs(items: KeyValuePair[]): KeyValuePair[] {
  const sorted = [...items];

  sorted.sort((a, b) => {
    if (a.key < b.key) {
      return -1;
    }
    if (a.key > b.key) {
      return 1;
    }
    return 0;
  });

  return sorted;
}

export default KeyValuePairTable;
