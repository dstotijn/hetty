import { KeyValuePair } from "./components/KeyValuePair";

function updateURLQueryParams(url: string, queryParams: KeyValuePair[]) {
  // Note: We don't use the `URL` interface, because we're potentially dealing
  // with malformed/incorrect URLs, which would yield TypeErrors when constructed
  // via `URL`.
  let newURL = url;

  const questionMarkIndex = url.indexOf("?");
  if (questionMarkIndex !== -1) {
    newURL = newURL.slice(0, questionMarkIndex);
  }

  const searchParams = new URLSearchParams();
  for (const { key, value } of queryParams.filter(({ key }) => key !== "")) {
    searchParams.append(key, value);
  }

  const rawQueryParams = decodeURI(searchParams.toString());

  if (rawQueryParams == "") {
    return newURL;
  }

  return newURL + "?" + rawQueryParams;
}

export default updateURLQueryParams;
