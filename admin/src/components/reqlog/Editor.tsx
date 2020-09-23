import dynamic from "next/dynamic";
const MonacoEditor = dynamic(import("react-monaco-editor"), { ssr: false });

const monacoOptions = {
  readOnly: true,
  wordWrap: "on",
  minimap: {
    enabled: false,
  },
};

type language = "html" | "typescript" | "json";

function editorDidMount() {
  return (window.MonacoEnvironment.getWorkerUrl = (moduleId, label) => {
    if (label === "json") return "/_next/static/json.worker.js";
    if (label === "html") return "/_next/static/html.worker.js";
    if (label === "javascript") return "/_next/static/ts.worker.js";

    return "/_next/static/editor.worker.js";
  });
}

function languageForContentType(contentType: string): language {
  switch (contentType) {
    case "text/html":
      return "html";
    case "application/json":
      return "json";
    case "application/javascript":
      return "typescript";
    default:
      return;
  }
}

interface Props {
  content: string;
  contentType: string;
}

function Editor({ content, contentType }: Props): JSX.Element {
  return (
    <MonacoEditor
      height={"600px"}
      language={languageForContentType(contentType)}
      theme="vs-dark"
      editorDidMount={editorDidMount}
      options={monacoOptions as any}
      value={content}
    />
  );
}

export default Editor;
