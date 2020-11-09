# Getting Started

## Installation

Hetty compiles to a static binary, with an embedded SQLite database and web
admin interface.

### Install pre-built release (recommended)

👉 Downloads for Linux, macOS and Windows are available on the [releases page](https://github.com/dstotijn/hetty/releases).

### Build from source

#### Prerequisites

- [Go](https://golang.org/)
- [Yarn](https://yarnpkg.com/)
- [go.rice](https://github.com/GeertJohan/go.rice)

Hetty depends on SQLite (via [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3))
and needs `cgo` to compile. Additionally, the static resources for the web admin interface
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
For persistent storage of CA certificate and project databases, mount a volume:

```
$ mkdir -p $HOME/.hetty
$ docker run -v $HOME/.hetty:/root/.hetty -p 8080:8080 dstotijn/hetty
```

## Usage

When Hetty is started, by default it listens on `:8080` and is accessible via
[http://localhost:8080](http://localhost:8080). Depending on incoming HTTP
requests, it either acts as a MITM proxy, or it serves the API and web interface.

By default, project database files and CA certificates are stored in a `.hetty`
directory under the user's home directory (`$HOME` on Linux/macOS, `%USERPROFILE%`
on Windows).

To start, ensure `hetty` (downloaded from a release, or manually built) is in your
`$PATH` and run:

```
$ hetty
```

You should see:

```
2020/11/01 14:47:10 [INFO] Running server on :8080 ...
```

Then, visit [http://localhost:8080](http://localhost:8080) to get started.

### Configuration

An overview of available configuration flags:

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
