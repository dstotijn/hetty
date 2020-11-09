---
sidebarDepth: 3
---

# Modules

Hetty consists of various _modules_ that together form an HTTP toolkit. They
typically are managed via the web admin interface. Some modules expose settings
and behavior that is leveraged by other modules.

The available modules:

[[toc]]

## Projects

Projects are self-contained (SQLite) database files that contain module data.
They allow you organize your work, for example to split your work between research
targets.

You can create multiple projects, but only one can be open at a time. Most other
modules are useful only if you have a project opened, so creating a project is
typically the first thing you do when you start using Hetty.

### Creating a new project

When you open the Hetty admin interface after starting the program, you’ll be prompted
on the homepage to create a new project. Give it a name (alphanumeric and space character)
and click the create button:

![Creating a project](./create_project.png =417x)

The project name will become the base for the database file on disk. For example,
if you name your project `My first project`, the file on disk will be
`My first project.db`.

::: tip INFO
Project database files by default are stored in `$HOME/.hetty/projects` on Linux
and macOS, and `%USERPROFILE%/.hetty` on Windows. You can override this path with
the `-projects` flag. See: [Usage](/guide/getting-started.md#usage).
:::

### Managing projects

You can open and delete existing projects on the “Projects” page, available via
the folder icon in the menu bar.

![Managing projects](./manage_projects.png =594x)

An opened (_active_) project is listed in green. You can close it using the “X”
button. To delete a project, use the trash bin icon.

::: danger
Deleting a project is irreversible.
:::

## Proxy

Hetty features a HTTP/1.1 proxy server with machine-in-the-middle (MITM) behavior.
For now, its only configuration is done via command line flags.

::: tip INFO
Support for HTTP/2 and WebSockets are currently not supported, but this will
likely be addressed in the (near) future.
:::

### Network address

To configure the network address that the proxy listens on, use the `-addr` flag
when starting Hetty. The address needs to be in the format `[host]:port`. E.g.
`localhost:3000` or `:3000`. If the host in the address is empty or a literal
unspecified IP address, Hetty listens on all available unicast and anycast IP
addresses of the local system.

::: tip INFO
When not specified with `-addr`, Hetty by default listens on `:8080`.
:::

Example of starting Hetty, binding to port `3000` on all IPs of the local system:

```
$ hetty -addr :3000
```

### Using the proxy

To use Hetty as an HTTP proxy server, you’ll need to configure your HTTP client (e.g.
your browser or mobile OS). Refer to your client documentation or use a search
engine to find instructions for setting a HTTP proxy.

### Certificate Authority (CA)

In order for Hetty to proxy requests going to HTTPS endpoints, a root CA certificate for
Hetty will need to be set up. Furthermore, the CA certificate needs to be
installed to the host for them to be trusted by your browser. The following steps
will cover how you can generate a certificate, provide it to Hetty, and how
you can install it in your local CA store.

::: tip INFO
Certificate management features (e.g. automated installing of a root CA to your local
OS or browser trust store) are planned for a future release. In the meantime, please
use the instructions below.
:::

#### Generating a CA certificate

You can generate a CA keypair two different ways. The first is bundled directly
with Hetty, and simplifies the process immensely. The alternative is using OpenSSL
to generate them, which provides more control over expiration time and cryptography
used, but requires you install the OpenSSL tooling. The first is suggested for any
beginners trying to get started.

#### Generating CA certificates with hetty

Hetty will generate the default key and certificate on its own if none are supplied
or found in `~/.hetty/` when first running the CLI. To generate a default key and
certificate with hetty, simply run the command with no arguments.

You should now have a key and certificate located at `~/.hetty/hetty_key.pem` and
`~/.hetty/hetty_cert.pem` respectively.

#### Generating CA certificates with OpenSSL

::: tip INFO
This following instructions are for Linux but should provide guidance for Windows
and macOS as well.
:::

You can start off by generating a new key and CA certificate which will both expire
after a month.

```

$ mkdir ~/.hetty
$ openssl req -newkey rsa:2048 -new -nodes -x509 -days 31 -keyout ~/.hetty/hetty_key.pem -out ~/.hetty/hetty_cert.pem

```

The default location which `hetty` will check for the key and CA certificate is under
`~/.hetty/`, at `hetty_key.pem` and `hetty_cert.pem` respectively. You can move them
here and `hetty` will detect them automatically. Otherwise, you can specify the
location of these as arguments to `hetty`.

```

\$ hetty -key /some/directory/key.pem -cert /some/directory/cert.pem

```

#### Trusting the CA certificate

In order for your browser to allow traffic to the local Hetty proxy, you may need
to install these certificates to your local CA store.

On Ubuntu, you can update your local CA store with the certificate by running the
following commands:

```

$ sudo cp ~/.hetty/hetty_cert.pem /usr/local/share/ca-certificates/hetty.crt
$ sudo update-ca-certificates

```

On Windows, you would add your certificate by using the Certificate Manager,
which you can run via:

```

certmgr.msc

```

On macOS, you can add your certificate by using the Keychain Access program. This
can be found under `Application/Utilities/Keychain Access.app`. After opening this,
drag the certificate into the app. Next, open the certificate in the app, enter the
_Trust_ section, and under _When using this certificate_ select _Always Trust_.

::: tip INFO
Various Linux distributions may require other steps or commands for updating
their certificate authority. See the documentation relevant to your distribution for
more information on how to update the system to trust your self-signed certificate.
:::

## Scope

The scope module lets you define _rules_ that other modules can use to control
their behavior. For example, the [proxy logs module](#proxy-logs) can be configured to only
show logs for in-scope requests; meaning only requests are shown that match one
or more scope rules.

### Managing scope rules

You can manage scope rules via the “Scope” page, available via the crosshair icon
in the menu bar.

A rule consists of a _type_ and a regular expression ([RE2 syntax](https://github.com/google/re2/wiki/Syntax)).
The only supported type at the moment is “URL”.

::: tip INFO
Just like all module configuration, scope rules are defined and stored per-project.
:::

#### Adding a rule

On the ”Scope” page, enter a regular expression and click “Add rule”:

![Adding a scope rule](./add_scope_rule.png =592x)

_Example: Rule that matches URLs with `example.com` (or any subdomain) on any path._

#### Deleting rules

Use the trash icon next to an existing scope rule to delete it.

## Proxy logs

You can few logs captured by the Proxy module on the Proxy logs page, available
via the proxy icon in the menu bar.

![Proxy logs overview](./proxy_logs.png =1207x)

### Showing a log entry

Click a row in the overview table to view log details in the bottom request and
response panes. When a request and/or response has a body, it's shown below the
HTTP headers. Header keys and values can be copied to clipboard by clicking them.

### Filtering logs

To only show log entries that match any of the [scope rules](#scope), click the
filter icon in the search bar and select “Only show in-scope requests”:

![Only show in-scope requests](./filter_in_scope.png =431x)

::: tip INFO
At the moment of writing (`v0.2.0`), text based search is not implemented yet.
:::
