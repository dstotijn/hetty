import { KeyValuePair } from "./components/KeyValuePair";

export function queryParamsFromURL(url: string): KeyValuePair[] {
  const questionMarkIndex = url.indexOf("?");
  if (questionMarkIndex === -1) {
    return [];
  }

  const queryParams: KeyValuePair[] = [];

  const searchParams = new URLSearchParams(url.slice(questionMarkIndex + 1));
  for (const [key, value] of searchParams) {
    queryParams.push({ key, value });
  }

  return queryParams;
}
