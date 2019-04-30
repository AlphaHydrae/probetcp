# tcpwait

Wait for TCP endpoints to be reachable (e.g. wait for a database in a Docker container).

```
$> tcpwait --retries 10 --timeout 1000 db.example.com:5432 -- rails server
```

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Installation](#installation)
  - [Dockerfile](#dockerfile)
  - [Homebrew](#homebrew)
  - [Download binary](#download-binary)
- [Usage](#usage)
  - [Timeout](#timeout)
  - [Retrying](#retrying)
  - [Multiple endpoints](#multiple-endpoints)
  - [Executing a command](#executing-a-command)
  - [Environment variable interpolation](#environment-variable-interpolation)
    - [Dockerfile](#dockerfile-1)
  - [Quiet](#quiet)
  - [TL;DR](#tldr)
- [Exit codes](#exit-codes)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

[![version](https://img.shields.io/badge/Version-v2.2.0-blue.svg)](https://github.com/AlphaHydrae/tcpwait/releases/tag/v2.2.0)
[![build status](https://travis-ci.org/AlphaHydrae/tcpwait.svg?branch=master)](https://travis-ci.org/AlphaHydrae/tcpwait)
[![license](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE.txt)



## Installation

### Dockerfile

```
RUN wget -O /usr/local/bin/tcpwait \
  https://github.com/AlphaHydrae/tcpwait/releases/download/v2.2.0/tcpwait_v2.2.0_linux_amd64 && \
  chmod +x /usr/local/bin/tcpwait
```

### Homebrew

```
brew install alphahydrae/tools/tcpwait
```

### Download binary

* **Linux**

  ```
  wget -O /usr/local/bin/tcpwait \
    https://github.com/AlphaHydrae/tcpwait/releases/download/v2.2.0/tcpwait_v2.2.0_linux_amd64 && \
    chmod +x /usr/local/bin/tcpwait
  ```
* **Linux (arm64)**

  ```
  wget -O /usr/local/bin/tcpwait \
    https://github.com/AlphaHydrae/tcpwait/releases/download/v2.2.0/tcpwait_v2.2.0_linux_arm64 && \
    chmod +x /usr/local/bin/tcpwait
  ```
* **macOS**

  ```
  wget -O /usr/local/bin/tcpwait \
    https://github.com/AlphaHydrae/tcpwait/releases/download/v2.2.0/tcpwait_v2.2.0_darwin_amd64 && \
    chmod +x /usr/local/bin/tcpwait
  ```
* **Windows**

  ```
  wget -O /usr/local/bin/tcpwait \
    https://github.com/AlphaHydrae/tcpwait/releases/download/v2.2.0/tcpwait_v2.2.0_windows_amd64 && \
    chmod +x /usr/local/bin/tcpwait
  ```



## Usage

**tcpwait** will attempt to establish a TCP connection to an endpoint and exit successfully if that endpoint can be reached:

```
$> tcpwait google.com:80
Reached "google.com:80" in 0.008598s

$> echo $?
0
```

If the connection cannot be established, **tcpwait** will exit with a non-zero code:

```
$> tcpwait google.com:12345
Error: could not reach "google.com:12345" after 1.000910s

$> echo $?
3
```

### Timeout

**tcpwait** uses a connection timeout of 1000 milliseconds by default.

Use the `-t, --timeout <value>` option to change it. For example, this command
will wait for one minute before determining that it cannot connect:

```
$> tcpwait -t 60000 google.com:12345
Error: could not reach "google.com:12345" after 60.004512s
```

Note that the timeout does not increase the duration of successful connection
attempts.  In this example, the command will return immediately even though the
timeout is high, because it stops as soon as the connection can be made:

```
$> tcpwait -t 60000 google.com:80
Reached "google.com:80" in 0.007431s
```

### Retrying

With the `-r, --retries <value>` option, **tcpwait** will retry to connect to
the TCP endpoint if it fails. The value of the option is the number of times to
retry:

```
$> tcpwait -r 4 google.com:12345
Waiting for google.com:12345 (1)...
Waiting for google.com:12345 (2)...
Waiting for google.com:12345 (3)...
Waiting for google.com:12345 (4)...
Error: could not reach "google.com:12345" after 5.019876s
```

> Note that the total number of attemps is equal to the number of retries plus
> one, in this example 5, because retries are made in addition to the initial
> connection attempt.

The `-i, --interval <value>` option can be used to introduce a delay between
each connection attempt:

```
$> tcpwait -i 500 -r 4 google.com:12345
Waiting for google.com:12345 (1)...
Waiting for google.com:12345 (2)...
Waiting for google.com:12345 (3)...
Waiting for google.com:12345 (4)...
Error: could not reach "google.com:12345" after 7.032964s
```

In this case, the total duration is 7 seconds instead of 5 like the previous
example, because the command waited for 500 additional milliseconds between each
connection attempt.

### Multiple endpoints

You may provide **tcpwait** with multiple TCP endpoints to attempt to connect
to. It will do so in parallel and return successfully as soon as all endpoints
have been successfully reached:

```
$> tcpwait github.com:22 github.com:80 github.com:443
Reached "github.com:443" in 0.030656s
Reached "github.com:22" in 0.030723s
Reached "github.com:80" in 0.030755s

$> echo $?
0
```

> Note that the order in which the result messages are printed is
> non-deterministic.

If one or more of the TCP endpoints cannot be reached, the command will exit
with a non-zero code:

```
$> tcpwait github.com:22 github.com:80 github.com:443 github.com:12345
Reached "github.com:80" in 0.029169s
Reached "github.com:443" in 0.029225s
Reached "github.com:22" in 0.029252s
Error: could not reach "github.com:12345" after 1.004060s

$> echo $?
3
```

You may use **tcpwait**'s retry and timeout options with multiple endpoints. All
connection attemts will be made with the same parameters. For example, this
command will attempt to reach `google.com:12345` and `google.com:23456` with 2
additional retry attempts before failing:

```
$> tcpwait -r 2 google.com:12345 google.com:23456
Waiting for google.com:23456 (1)...
Waiting for google.com:12345 (1)...
Waiting for google.com:23456 (2)...
Waiting for google.com:12345 (2)...
Error: could not reach "google.com:23456" after 3.013016s
```

### Executing a command

**tcpwait** will interpret all arguments after the terminator (`--`) as a
command to execute after all the TCP endpoints have been reached.

This can be used to conditionally execute a command as soon as a service, such
as a database, is reachable. For example, the following command could be used to
run a Ruby on Rails application as soon as the database server can be reached:

    tcpwait db.example.com:5432 -- rails server

> Note that the command is executed with [execve] and replaces **tcpwait**'s
> process, i.e. there is no leftover **tcpwait** process once the command
> executes, even if it is long-running like a web application.

### Environment variable interpolation

**tcpwait** will interpolate variables in TCP endpoint strings as well as in the
value of its numeric options (`-i, --interval`, `-r, --retries`, and `-t,
--timeout`):

```
$> export TCPWAIT_HOST=google.com TCPWAIT_PORT=12345

$> tcpwait --retries '${TCPWAIT_RETRIES-3}' '$TCPWAIT_HOST:$TCPWAIT_PORT'
Waiting for google.com:12345 (1)...
Waiting for google.com:12345 (2)...
Waiting for google.com:12345 (3)...
Error: could not reach "google.com:12345" after 4.011786s
```

Supported expansions are documented in the [interpolate] library.

> Note the use of single quotes. If you used double quotes or no quotes, the
> shell would interpolate the variables before they are passed to **tcpwait**.
> With single quotes, the string value, e.g. `${TCPWAIT_RETRIES-3}` is passed as
> is to **tcpwait**, which does the interpolation itself.

#### Dockerfile

An example of how this feature can be useful is to support interpolating
variables into **tcpwait**'s arguments when using it with the exec form of a
Dockerfile's [`ENTRYPOINT`][entrypoint]:

```
ENTRYPOINT [ "tcpwait", "${DB_HOST-db}:${DB_PORT-5432}", "--", "rails", "server" ]
```

Normally these variables would only be interpolated using `ENTRYPOINT`'s *shell*
form, but since **tcpwait** is capable of interpolating those variables at
runtime, you can use it with `ENTRYPOINT`'s *exec* form.

### Quiet

**tcpwait** outputs information on its standard error stream by default, to
indicate failure or success, or that it it waiting. You can silence these
messages with the `-q, --quiet` option:

```
$> tcpwait -q -r 4 google.com:12345

$> echo $?
3
```

### TL;DR

Run `tcpwait --help` for instructions:

```
tcpwait waits for TCP endpoints to be reachable.

Usage:
  tcpwait [OPTION...] ENDPOINT... [--] [EXEC...]

Options:
  -i, --interval string   Time to wait between retries in milliseconds (default 0) (default "0")
  -q, --quiet             Do not print anything (default false)
  -r, --retries string    Number of times to retry to reach the endpoint if it fails (default 0) (default "0")
  -t, --timeout string    TCP connection timeout in milliseconds (default 1000) (default "1000")

Examples:
  Wait for a website:
    tcpwait google.com:80
  Wait for a MySQL database (10 attempts every 2 seconds):
    tcpwait -r 9 -i 2000 tcp://localhost:3306
  Wait for multiple endpoints:
    tcpwait github.com:22 github.com:80 github.com:443
  Execute a command after an endpoint is reached:
    tcpwait db.example.com:5432 -- rails server
```



## Exit codes

**tcpwait** may exit with the following status codes:

Code | Description
:--- | :---
`0`  | All endpoints were reached successfully.
`1`  | Invalid arguments were given.
`2`  | An unrecoverable error occurred while attempting to reach a TCP endpoint.
`3`  | One of the endpoints could not be reached (even after retrying, if applicable).
`4`  | An unrecoverable error occurred while trying to interpolate environment variables into options or TCP endpoint strings.
`10` | The command to execute (provided after `--`) could not be found in the `$PATH`.
`11` | An unrecoverable error occurred while attempting to execute the command.

> Note that if a command to execute is specified (after `--`), it is executed
> with [execve], meaning that the **tcpwait** process is replaced by the
> command's.
>
> In this case, the exit code returned will be that of the executed command, not
> **tcpwait**'s. Look in the command's documentation for the meaning of its exit
> codes.



[entrypoint]: https://docs.docker.com/engine/reference/builder/#entrypoint
[execve]: https://linux.die.net/man/2/execve
[interpolate]: https://github.com/buildkite/interpolate#supported-expansions
[go]: https://golang.org
