function omitTypename<T>(key: string, value: T): T | undefined {
  return key === "__typename" ? undefined : value;
}

export function withoutTypename<T>(input: T): T {
  return JSON.parse(JSON.stringify(input), omitTypename);
}
