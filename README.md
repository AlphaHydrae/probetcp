# tcpwait

Wait for TCP endpoints to be reachable (e.g. wait for a database in a Docker container).

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Usage](#usage)
- [Installation](#installation)
  - [Download binary](#download-binary)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

![version](https://img.shields.io/badge/Version-v2.0.0-blue.svg)
[![build status](https://travis-ci.org/AlphaHydrae/tcpwait.svg?branch=master)](https://travis-ci.org/AlphaHydrae/tcpwait)
[![license](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE.txt)



## Usage

```
tcpwait waits for TCP endpoints to be reachable.

Usage:
  tcpwait [OPTION...] ENDPOINT... [--] [EXEC...]

Options:
  -i, --interval int   Time to wait between retries in milliseconds (default 1000)
  -q, --quiet          Do not print anything (default false)
  -r, --retries int    Number of times to retry to reach the endpoint if it fails (default 0)
  -t, --timeout int    TCP connection timeout in milliseconds (default 60000)

Examples:
  Wait for a website:
    tcpwait google.com:80
  Wait for a MySQL database (10 attempts every 2 seconds):
    tcpwait -r 9 -i 2000 tcp://localhost:3306
  Wait for multiple endpoints:
    tcpwait github.com:22 github.com:80 github.com:443
  Execute a command after an endpoint is reached:
    tcpwait db.example.com:5432 -- pg_dump -H db.example.com -f dump.sql example
```



## Installation

With [Homebrew][brew]:

```
brew install alphahydrae/tools/tcpwait
```

### Download binary

* **Dockerfile**

  ```
  RUN wget -O /usr/local/bin/tcpwait \
    https://github.com/AlphaHydrae/tcpwait/releases/download/v2.0.0/tcpwait_v2.0.0_linux_amd64 && \
    chmod +x /usr/local/bin/tcpwait
  ```
* **Linux**

  ```
  wget -O /usr/local/bin/tcpwait \
    https://github.com/AlphaHydrae/tcpwait/releases/download/v2.0.0/tcpwait_v2.0.0_linux_amd64 && \
    chmod +x /usr/local/bin/tcpwait
  ```
* **Linux (arm64)**

  ```
  wget -O /usr/local/bin/tcpwait \
    https://github.com/AlphaHydrae/tcpwait/releases/download/v2.0.0/tcpwait_v2.0.0_linux_arm64 && \
    chmod +x /usr/local/bin/tcpwait
  ```
* **macOS**

  ```
  wget -O /usr/local/bin/tcpwait \
    https://github.com/AlphaHydrae/tcpwait/releases/download/v2.0.0/tcpwait_v2.0.0_darwin_amd64 && \
    chmod +x /usr/local/bin/tcpwait
  ```
* **Windows**

  ```
  wget -O /usr/local/bin/tcpwait \
    https://github.com/AlphaHydrae/tcpwait/releases/download/v2.0.0/tcpwait_v2.0.0_windows_amd64 && \
    chmod +x /usr/local/bin/tcpwait
  ```



## Exit codes

**tcpwait** may exit with the following status codes:

Code | Description
:--- | :---
`0`  | All endpoints were reached successfully.
`1`  | Invalid arguments were given.
`2`  | An unrecoverable error occurred while attempting to reach a TCP endpoint.
`3`  | One of the endpoints could not be reached (even after retrying, if applicable).
`10` | The command to execute (provided after `--`) could not be found in the `$PATH`.
`11` | An unrecoverable error occurred while attempting to execute the command.



[brew]: https://brew.sh/
[go]: https://golang.org
