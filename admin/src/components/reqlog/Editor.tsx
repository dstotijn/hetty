import MonacoEditor from "@monaco-editor/react";
import monaco from "monaco-editor/esm/vs/editor/editor.api";

const monacoOptions: monaco.editor.IEditorOptions = {
  readOnly: true,
  wordWrap: "on",
  minimap: {
    enabled: false,
  },
};

type language = "html" | "typescript" | "json";

function languageForContentType(contentType?: string): language | undefined {
  switch (contentType) {
    case "text/html":
      return "html";
    case "application/json":
    case "application/json; charset=utf-8":
      return "json";
    case "application/javascript":
    case "application/javascript; charset=utf-8":
      return "typescript";
    default:
      return;
  }
}

interface Props {
  content: string;
  contentType?: string;
}

function Editor({ content, contentType }: Props): JSX.Element {
  return (
    <MonacoEditor
      height={"600px"}
      language={languageForContentType(contentType)}
      theme="vs-dark"
      options={monacoOptions}
      value={content}
    />
  );
}

export default Editor;
