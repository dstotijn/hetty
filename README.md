# Ratatoskr

Ratatoskr is a free, open-source web security proxy for penetration testers and
bug bounty hunters. It ships as a single binary with a Go backend, a Svelte web
UI, SQLite project storage and a GraphQL API.

> Status: early development (milestone M0). Nothing is usable yet.

## Build

Requires Go 1.27 and, once the web UI lands, [Bun](https://bun.sh).

```sh
make build   # ./ratatoskr
make check   # go fix -diff, go vet, golangci-lint, go test -race, web checks
make run     # build and start with ./.dev-data
```

## Documentation

- [Project overview](docs/00-overview.md)
- [Feature spec](docs/01-feature-spec.md)
- [Architecture](docs/03-architecture.md)
- [Roadmap](docs/08-roadmap.md)
- [Contributor and agent rules](AGENTS.md)

## Credits

Ratatoskr is a fork of [Hetty](https://github.com/dstotijn/hetty) by David
Stotijn. Hetty's proxy, certificate handling, intercept queue and API gave this
project its starting point.

## License

[MIT](LICENSE).
