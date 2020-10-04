<img src="https://i.imgur.com/AT71SBq.png" width="346" />

> Hetty is an HTTP toolkit for security research. It aims to become an open source
> alternative to commercial software like Burp Suite Pro, with powerful features
> tailored to the needs of the infosec and bug bounty community.

<img src="https://i.imgur.com/ZZ6o83X.png">

## Features/to do

- [x] HTTP man-in-the-middle (MITM) proxy and GraphQL server.
- [x] Web interface (Next.js) with proxy log viewer.
- [ ] Add scope support to the proxy.
- [ ] Full text search (with regex) in proxy log viewer.
- [ ] Project management.
- [ ] Sender module for sending manual HTTP requests, either from scratch or based
      off requests from the proxy log.
- [ ] Attacker module for automated sending of HTTP requests. Leverage the concurrency
      features of Go and its `net/http` package to make it blazingly fast.

## Installation

Hetty is packaged on GitHub as a single binary, with the web interface resources
embedded.

👉 You can find downloads for Linux, macOS and Windows on the [releases page](https://github.com/dstotijn/hetty/releases).

### Alternatives:

**Build from source**

```
$ GO111MODULE=auto go get -u -v github.com/dstotijn/hetty/cmd/hetty
```

Then export the Next.js frontend app:

```
$ cd admin
$ yarn install
$ yarn export
```

This will ensure a folder `./admin/dist` exists.
Then, you can bundle the frontend app using `rice`.
The easiest way to do this is via a supplied `Makefile` command in the root of
the project:

```
make build
```

**Docker**

Alternatively, you can run Hetty via Docker. See: [`dstotijn/hetty`](https://hub.docker.com/r/dstotijn/hetty)
on Docker Hub.

```
$ docker run \
-v $HOME/.hetty/hetty_key.pem:/root/.hetty/hetty_key.pem \
-v $HOME/.hetty/hetty_cert.pem:/root/.hetty/hetty_cert.pem \
-v $HOME/.hetty/hetty.db:/root/.hetty/hetty.db \
-p 127.0.0.1:8080:8080 \
dstotijn/hetty
```

## Usage

Hetty is packaged as a single binary, with the web interface resources embedded.
When the program is run, it listens by default on `:8080` and is accessible via
http://localhost:8080. Depending on incoming HTTP requests, it either acts as a
MITM proxy, or it serves the GraphQL API and web interface (Next.js).

```
$ hetty -h
Usage of ./hetty:
  -addr string
        TCP address to listen on, in the form "host:port" (default ":8080")
  -adminPath string
        File path to admin build
  -cert string
        CA certificate filepath. Creates a new CA certificate is file doesn't exist (default "~/.hetty/hetty_cert.pem")
  -db string
        Database file path (default "~/.hetty/hetty.db")
  -key string
        CA private key filepath. Creates a new CA private key if file doesn't exist (default "~/.hetty/hetty_key.pem")
```

⚠️ _Todo: Write instructions for installing CA certificate in local CA store, and_
_configuring Hetty to be used as a proxy server._

## Vision and roadmap

The project has just gotten underway, and as such I haven’t had time yet to do a
write-up on its mission and roadmap. A short summary/braindump:

- Fast core/engine, built with Go, with a minimal memory footprint.
- GraphQL server to interact with the backend.
- Easy to use web interface, built with Next.js and Material UI.
- Extensibility is top of mind. All modules are written as Go packages, to
  be used by the main `hetty` program, but also usable as libraries for other software.
  Aside from the GraphQL server, it should (eventually) be possible to also use
  it as a CLI tool.
- Pluggable architecture for the MITM proxy and future modules, making it
  possible for hook into the core engine.
- Talk to the community, and focus on the features that the majority.
  Less features means less code to maintain.

## Status

The project is currently under active development. Please star/follow and check
back soon. 🤗

## Contributing

Please see the [Contribution Guidelines](CONTRIBUTING.md) for details.

## Acknowledgements

Thanks to the [Hacker101 community on Discord](https://www.hacker101.com/discord)
for all the encouragement to actually start building this thing!

## License

[MIT](LICENSE)

---

© 2020 David Stotijn — [Twitter](https://twitter.com/dstotijn), [Email](mailto:dstotijn@gmail.com)
