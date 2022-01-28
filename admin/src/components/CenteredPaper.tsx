import { Paper } from "@mui/material";

function CenteredPaper({ children }: { children: React.ReactNode }): JSX.Element {
  return (
    <div>
      <Paper
        elevation={0}
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          padding: 36,
        }}
      >
        {children}
      </Paper>
    </div>
  );
}

export default CenteredPaper;
