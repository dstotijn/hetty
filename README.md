<svg width="331" height="73" viewBox="0 0 331 73" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M0.4 60V8.16H9.4V28.896H23.728V8.16H32.728V60H23.728V37.176H9.4V60H0.4ZM42.5639 35.16C42.5639 32.808 42.9719 30.672 43.7879 28.752C44.6039 26.832 45.7559 25.2 47.2439 23.856C48.7319 22.512 50.5319 21.48 52.6439 20.76C54.7559 20.04 57.1319 19.68 59.7719 19.68C62.3639 19.68 64.6919 20.064 66.7559 20.832C68.8679 21.552 70.6679 22.584 72.1559 23.928C73.6439 25.272 74.7959 26.904 75.6119 28.824C76.4279 30.744 76.8359 32.88 76.8359 35.232V42.432H51.4919V45.024C51.4919 47.712 52.2119 49.8 53.6519 51.288C55.1399 52.776 57.1799 53.52 59.7719 53.52C61.6919 53.52 63.3479 53.208 64.7399 52.584C66.1319 51.96 67.1159 51.096 67.6919 49.992H76.5479C76.1159 51.624 75.3959 53.112 74.3879 54.456C73.3799 55.752 72.1319 56.88 70.6439 57.84C69.2039 58.752 67.5479 59.472 65.6759 60C63.8519 60.48 61.8839 60.72 59.7719 60.72C57.1799 60.72 54.8039 60.36 52.6439 59.64C50.5319 58.92 48.7319 57.888 47.2439 56.544C45.7559 55.152 44.6039 53.496 43.7879 51.576C42.9719 49.656 42.5639 47.52 42.5639 45.168V35.16ZM51.4919 36.384H68.0519V35.232C68.0519 32.592 67.3079 30.552 65.8199 29.112C64.3799 27.624 62.3639 26.88 59.7719 26.88C57.1799 26.88 55.1399 27.624 53.6519 29.112C52.2119 30.552 51.4919 32.568 51.4919 35.16V36.384ZM84.5838 20.4H95.5278V8.16H104.528V20.4H119.576V28.536H104.528V47.76C104.528 49.008 104.888 50.016 105.608 50.784C106.376 51.504 107.432 51.864 108.776 51.864H118.856V60H108.128C104.24 60 101.168 58.896 98.9118 56.688C96.6558 54.48 95.5278 51.504 95.5278 47.76V28.536H84.5838V20.4ZM127.756 20.4H138.7V8.16H147.7V20.4H162.748V28.536H147.7V47.76C147.7 49.008 148.06 50.016 148.78 50.784C149.548 51.504 150.604 51.864 151.948 51.864H162.028V60H151.3C147.412 60 144.34 58.896 142.084 56.688C139.828 54.48 138.7 51.504 138.7 47.76V28.536H127.756V20.4ZM170.208 20.4H180.072L187.848 41.496C188.184 42.36 188.448 43.272 188.64 44.232C188.88 45.144 189.048 46.008 189.144 46.824C189.288 47.736 189.408 48.648 189.504 49.56H190.08C190.128 48.696 190.224 47.784 190.368 46.824C190.464 46.008 190.608 45.144 190.8 44.232C190.992 43.272 191.232 42.36 191.52 41.496L198.864 20.4H208.367L189.504 72.24H180L185.472 57.768L170.208 20.4Z" fill="black"/>
<path d="M225.835 54.312C225.835 52.296 226.435 50.688 227.635 49.488C228.883 48.24 230.491 47.616 232.459 47.616C234.427 47.616 236.035 48.24 237.283 49.488C238.531 50.688 239.155 52.296 239.155 54.312C239.155 56.184 238.531 57.72 237.283 58.92C236.035 60.12 234.427 60.72 232.459 60.72C230.491 60.72 228.883 60.12 227.635 58.92C226.435 57.72 225.835 56.184 225.835 54.312ZM225.835 26.376C225.835 24.36 226.435 22.752 227.635 21.552C228.883 20.304 230.491 19.68 232.459 19.68C234.427 19.68 236.035 20.304 237.283 21.552C238.531 22.752 239.155 24.36 239.155 26.376C239.155 28.248 238.531 29.784 237.283 30.984C236.035 32.184 234.427 32.784 232.459 32.784C230.491 32.784 228.883 32.184 227.635 30.984C226.435 29.784 225.835 28.248 225.835 26.376ZM321.323 0.239997H330.683L305.555 67.92H296.123L321.323 0.239997ZM288.923 0.239997H298.283L273.155 67.92H263.723L288.923 0.239997Z" fill="#1DE9B6"/>
</svg>

> Hetty is an HTTP toolkit for security research. It aims to become an open source
> alternative to commercial software like Burp Suite Pro, with powerful features
> tailored to the needs of the infosec and bug bounty community.

<img src="https://i.imgur.com/ZZ6o83X.png">

## Features/to do

- [x] HTTP man-in-the-middle (MITM) proxy and GraphQL server.
- [x] Web interface (Next.js) with proxy log viewer.
- [] Add scope support to the proxy.
- [] Full text search (with regex) in proxy log viewer.
- [] Project management.
- [] Sender module for sending manual HTTP requests, either from scratch or based
  off requests from the proxy log.
- [] Attacker module for automated sending of HTTP requests. Leverage the concurrency
  features of Go and its `net/http` package to make it blazingly fast.

## Installation

Hetty is packaged on GitHub as a single binary, with the web interface resources
embedded.

👉 You can find downloads for Linux, macOS and Windows on the [releases page](https://github.com/dstotijn/hetty/releases).

### Alternatives:

**Build from source**

```
$ go get github.com/dstotijn/hetty
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
-v $HOME/.ssh/hetty_key.pem:/.ssh/hetty_key.pem \
-v $HOME/.ssh/hetty_cert.pem:/.ssh/hetty_cert.pem \
-v $HOME/.hetty/hetty.db:/app/hetty.db \
-p 127.0.0.1:8080:80 \
dstotijn/hetty -key /.ssh/hetty_key.pem -cert /.ssh/hetty_cert.pem -db hetty.db
```

## Usage

Hetty is packaged as a single binary, with the web interface resources embedded.
When the program is run, it listens by default on `:8080` and is accessible via
http://localhost:8080. Depending on incoming HTTP requests, it either acts as a
MITM proxy, or it serves the GraphQL API and web interface (Next.js).

```
$ hetty -h
Usage of hetty:
  -addr string
    	TCP address to listen on, in the form "host:port" (default ":80")
  -adminPath string
    	File path to admin build
  -cert string
    	CA certificate file path
  -db string
    	Database file path (default "hetty.db")
  -key string
    	CA private key file path
```

**Note:** There is no built-in in support yet for generating a CA certificate.
This will be added really soon in an upcoming release. In the meantime, please
use `openssl` (_TODO: add instructions_).

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
- I’ve chosen [Cayley](https://cayley.io/) as the graph database (backed by
  BoltDB storage on disk) for now (not sure if it will work in the long run).
  The benefit is that Cayley (also written in Go)
  is embedded as a library. Because of this, the complete application is self contained
  in a single running binary.
- Talk to the community, and focus on the features that the majority.
  Less features means less code to maintain.

## Status

The project is currently under active development. Please star/follow and check
back soon. 🤗

## Acknowledgements

Thanks to the [Hacker101 community on Discord](https://discordapp.com/channels/514337135491416065)
for all the encouragement to actually start building this thing!

## License

[MIT](LICENSE)

---

© 2020 David Stotijn — [Twitter](https://twitter.com/dstotijn), [Email](mailto:dstotijn@gmail.com)
