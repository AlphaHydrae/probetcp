# probetcp

Probe TCP endpoints.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Installation](#installation)
- [Usage](#usage)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation

* **Linux**

  ```
  wget -O /usr/local/bin/probetcp \
    https://github.com/AlphaHydrae/probetcp/releases/download/v1.0.0/probetcp_v1.0.0_linux_amd64 && \
    chmod +x /usr/local/bin/probetcp
  ```
* **Linux (arm64)**

  ```
  wget -O /usr/local/bin/probetcp \
    https://github.com/AlphaHydrae/probetcp/releases/download/v1.0.0/probetcp_v1.0.0_linux_arm64 && \
    chmod +x /usr/local/bin/probetcp
  ```
* **macOS**

  ```
  wget -O /usr/local/bin/probetcp \
    https://github.com/AlphaHydrae/probetcp/releases/download/v1.0.0/probetcp_v1.0.0_darwin_amd64 && \
    chmod +x /usr/local/bin/probetcp
  ```
* **Windows**

  ```
  wget -O /usr/local/bin/probetcp \
    https://github.com/AlphaHydrae/probetcp/releases/download/v1.0.0/probetcp_v1.0.0_windows_amd64 && \
    chmod +x /usr/local/bin/probetcp
  ```

## Usage

```
probetcp probes TCP endpoints.

Usage:
  probetcp [OPTION...] TARGET

Options:
  -i, --interval int   Time to wait between probe retries in milliseconds (default 1000)
  -q, --quiet          Do not print anything (default false)
  -r, --retries int    Number of times to retry to probe the target if it fails (default 0)
  -t, --timeout int    TCP connection timeout in milliseconds (default 60000)

Examples:
  Probe a MySQL database over TCP (10 attempts every 2 seconds):
    probecli -r 9 -i 2000 tcp://localhost:3306
```
