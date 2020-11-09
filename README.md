<h1>
  <a href="https://github.com/dstotijn/hetty">
    <img src="https://hetty.xyz/assets/logo.png" width="293">
  </a>
</h1>

[![Latest GitHub release](https://img.shields.io/github/v/release/dstotijn/hetty?color=18BA91&style=flat-square)](https://github.com/dstotijn/hetty/releases/latest)
![GitHub download count](https://img.shields.io/github/downloads/dstotijn/hetty/total?color=18BA91&style=flat-square)
[![GitHub](https://img.shields.io/github/license/dstotijn/hetty?color=18BA91&style=flat-square)](https://github.com/dstotijn/hetty/blob/master/LICENSE)
[![Documentation](https://img.shields.io/badge/hetty-docs-18BA91?style=flat-square)](https://hetty.xyz/)

**Hetty** is an HTTP toolkit for security research. It aims to become an open
source alternative to commercial software like Burp Suite Pro, with powerful
features tailored to the needs of the infosec and bug bounty community.

<img src="https://hetty.xyz/assets/hetty_v0.2.0_header.png">

## Features

- Man-in-the-middle (MITM) HTTP/1.1 proxy with logs
- Project based database storage (SQLite)
- Scope support
- Headless management API using GraphQL
- Embedded web interface (Next.js)

ℹ️ Hetty is in early development. Additional features are planned
for a `v1.0` release. Please see the <a href="https://github.com/dstotijn/hetty/projects/1">backlog</a>
for details.

## Documentation

📖 [Read the docs.](https://hetty.xyz/)

## Installation

Hetty compiles to a self-contained binary, with an embedded SQLite database
and web based admin interface.

### Install pre-built release (recommended)

👉 Downloads for Linux, macOS and Windows are available on the [releases page](https://github.com/dstotijn/hetty/releases).

### Build from source

#### Prerequisites

- [Go](https://golang.org/)
- [Yarn](https://yarnpkg.com/)
- [go.rice](https://github.com/GeertJohan/go.rice)

Hetty depends on SQLite (via [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3))
and needs `cgo` to compile. Additionally, the static resources for the admin interface
(Next.js) need to be generated via [Yarn](https://yarnpkg.com/) and embedded in
a `.go` file with [go.rice](https://github.com/GeertJohan/go.rice) beforehand.

Clone the repository and use the `build` make target to create a binary:

```
$ git clone git@github.com:dstotijn/hetty.git
$ cd hetty
$ make build
```

### Docker

A Docker image is available on Docker Hub: [`dstotijn/hetty`](https://hub.docker.com/r/dstotijn/hetty).
For persistent storage of CA certificates and project databases, mount a volume:

```
$ mkdir -p $HOME/.hetty
$ docker run -v $HOME/.hetty:/root/.hetty -p 8080:8080 dstotijn/hetty
```

## Usage

When Hetty is run, by default it listens on `:8080` and is accessible via
http://localhost:8080. Depending on incoming HTTP requests, it either acts as a
MITM proxy, or it serves the API and web interface.

By default, project database files and CA certificates are stored in a `.hetty`
directory under the user's home directory (`$HOME` on Linux/macOS, `%USERPROFILE%`
on Windows).

To start, ensure `hetty` (downloaded from a release, or manually built) is in your
`$PATH` and run:

```
$ hetty
```

An overview of configuration flags:

```
$ hetty -h
Usage of ./hetty:
  -addr string
        TCP address to listen on, in the form "host:port" (default ":8080")
  -adminPath string
        File path to admin build
  -cert string
        CA certificate filepath. Creates a new CA certificate is file doesn't exist (default "~/.hetty/hetty_cert.pem")
  -key string
        CA private key filepath. Creates a new CA private key if file doesn't exist (default "~/.hetty/hetty_key.pem")
  -projects string
        Projects directory path (default "~/.hetty/projects")
```

You should see:

```
2020/11/01 14:47:10 [INFO] Running server on :8080 ...
```

Then, visit [http://localhost:8080](http://localhost:8080) to get started.

ℹ️ Detailed documentation is under development and will be available soon.

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

You should now have a key and certificate located at `~/.hetty/hetty_key.pem` and
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
to install these certificates to your local CA store.

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

- Fast core/engine, built with Go, with a minimal memory footprint.
- Easy to use admin interface, built with Next.js and Material UI.
- Headless management, via GraphQL API.
- Extensibility is top of mind. All modules are written as Go packages, to
  be used by Hetty, but also as libraries by other software.
- Pluggable architecture for MITM proxy, projects, scope. It should be possible.
  to build a plugin system in the (near) future.
- Based on feedback and real-world usage of pentesters and bug bounty hunters.
- Aim for a relatively small core feature set that the majority of security researchers need.

## Support

Use [issues](https://github.com/dstotijn/hetty/issues) for bug reports and
feature requests, and [discussions](https://github.com/dstotijn/hetty/discussions)
for questions and troubleshooting.

## Community

💬 [Join the Hetty Discord server](https://discord.gg/3HVsj5pTFP).

## Contributing

Want to contribute? Great! Please check the [Contribution Guidelines](CONTRIBUTING.md)
for details.

## Acknowledgements

- Thanks to the [Hacker101 community on Discord](https://www.hacker101.com/discord)
  for all the encouragement and feedback.
- The font used in the logo and admin interface is [JetBrains Mono](https://www.jetbrains.com/lp/mono/).

## License

[MIT License](LICENSE)

---

© 2020 David Stotijn — [Twitter](https://twitter.com/dstotijn), [Email](mailto:dstotijn@gmail.com)
