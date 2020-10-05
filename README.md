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
-v $HOME/.hetty/hetty.bolt:/root/.hetty/hetty.bolt \
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
        Database file path (default "~/.hetty/hetty.bolt")
  -key string
        CA private key filepath. Creates a new CA private key if file doesn't exist (default "~/.hetty/hetty_key.pem")
```

## Certificate Setup and Installation

In order for Hetty to proxy requests going to HTTPS endpoints, a root CA certificate for 
Hetty will need to be set up. Furthermore, the CA certificate may need to be
installed to the host for them to be trusted by your browser. The following steps
will cover how you can generate your certificate, provide them to hetty, and how
you can install them in your local CA store.

⚠️ _This process was done on a Linux machine but should_
_provide guidance on Windows and macOS as well._

### Generating CA certificates

You can generate a CA keypair two different ways. The first is bundled directly
with Hetty, and simplifies the process immensely. The alternative is using OpenSSL
to generate them, which provides more control over expiration time and cryptography
used, but requires you install the OpenSSL tooling. The first is suggested for any
beginners trying to get started.

#### Generating CA certificates with hetty

Hetty will generate the default key and certificate on its own if none are supplied
or found in `~/.hetty/` when first running the CLI. To generate a default key and
certificate with hetty, simply run the command with no arguments

```sh
hetty
```

You should now have a key and certificate located at  `~/.hetty/hetty_key.pem` and
`~/.hetty/hetty_cert.pem` respectively.

#### Generating CA certificates with OpenSSL

You can start off by generating a new key and CA certificate which will both expire
after a month.

```sh
mkdir ~/.hetty
openssl req -newkey rsa:2048 -new -nodes -x509 -days 31 -keyout ~/.hetty/hetty_key.pem -out ~/.hetty/hetty_cert.pem
```

The default location which `hetty` will check for the key and CA certificate is under
`~/.hetty/`, at `hetty_key.pem` and `hetty_cert.pem` respectively. You can move them
here and `hetty` will detect them automatically. Otherwise, you can specify the
location of these as arguments to `hetty`.

```
hetty -key key.pem -cert cert.pem
```

### Trusting the CA certificate

In order for your browser to allow traffic to the local Hetty proxy, you may need
to install these certificates to your local CA store. Otherwise, you may be met
with errors related to HTTP Strict Transport Security (HSTS). In the real world this
would be protecting you from man-in-the-middle attacks, but here we **are** the man.
Let's fix that.

On Ubuntu, you can update your local CA store with the certificate by running the
following commands:

```sh
sudo cp ~/.hetty/hetty_cert.pem /usr/local/share/ca-certificates/hetty.crt
sudo update-ca-certificates
```

On Windows, you would add your certificate by using the Certificate Manager. You
can launch that by running the command:

```batch
certmgr.msc
```

On macOS, you can add your certificate by using the Keychain Access program. This
can be found under `Application/Utilities/Keychain Access.app`. After opening this,
drag the certificate into the app. Next, open the certificate in the app, enter the
_Trust_ section, and under _When using this certificate_ select _Always Trust_.

_Note: Various Linux distributions may require other steps or commands for updating_
_their certificate authority. See the documentation relevant to your distribution for_
_more information on how to update the system to trust your self-signed certificate._

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

Thanks to the [Hacker101 community on Discord](https://www.hacker101.com/discord)
for all the encouragement to actually start building this thing!

## License

[MIT](LICENSE)

---

© 2020 David Stotijn — [Twitter](https://twitter.com/dstotijn), [Email](mailto:dstotijn@gmail.com)
