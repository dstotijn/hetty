const omitTypename = (key, value) => (key === "__typename" ? undefined : value);

export function withoutTypename(input: any): any {
  return JSON.parse(JSON.stringify(input), omitTypename);
}
