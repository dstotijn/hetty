import { KeyValuePair } from "./components/KeyValuePair";

function updateKeyPairItem(key: string, value: string, idx: number, items: KeyValuePair[]): KeyValuePair[] {
  const updated = [...items];
  updated[idx] = { key, value };

  // Append an empty key-value pair if the last item in the array isn't blank
  // anymore.
  if (items.length - 1 === idx && items[idx].key === "" && items[idx].value === "") {
    updated.push({ key: "", value: "" });
  }

  return updated;
}

export default updateKeyPairItem;
