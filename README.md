# probetcp

Probe TCP endpoints.

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
