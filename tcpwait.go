// The tcpwait command checks whether a TCP endpoint can be reached.
package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/alphahydrae/tcpwait/tcp"
	"github.com/fatih/color"
	flag "github.com/spf13/pflag"
)

const usageHeader = `%s waits for TCP endpoints to be reachable.

Usage:
  %s [OPTION...] ENDPOINT

Options:
`

const usageFooter = `
Examples:
  Wait for a website:
    tcpwait google.com:80
  Wait for a MySQL database (10 attempts every 2 seconds):
    tcpwait -r 9 -i 2000 tcp://localhost:3306
`

const tcpTargetRegexp = "^(?:tcp:\\/\\/)?(.+)$"

func main() {

	var interval int
	var retries int
	var quiet bool
	var timeout int

	flag.CommandLine.SetOutput(os.Stdout)

	flag.IntVarP(&interval, "interval", "i", 1e3, "Time to wait between retries in milliseconds")
	flag.IntVarP(&retries, "retries", "r", 0, "Number of times to retry to reach the endpoint if it fails (default 0)")
	flag.BoolVarP(&quiet, "quiet", "q", false, "Do not print anything (default false)")
	flag.IntVarP(&timeout, "timeout", "t", 60e3, "TCP connection timeout in milliseconds")

	flag.Usage = func() {
		fmt.Printf(usageHeader, os.Args[0], os.Args[0])
		flag.PrintDefaults()
		fmt.Print(usageFooter)
	}

	flag.Parse()

	endpoint := flag.Arg(0)

	if interval < 0 {
		fail(quiet, "the \"interval\" option must be greater than or equal to zero")
	} else if retries < 0 {
		fail(quiet, "the \"retries\" option must be greater than or equal to zero")
	} else if timeout <= 0 {
		fail(quiet, "the \"timeout\" option must be greater than zero")
	} else if endpoint == "" {
		fail(quiet, "an endpoint to wait for must be given as an argument (e.g. \"tcp://localhost:3306\")")
	}

	tcpRegexp := regexp.MustCompile(tcpTargetRegexp)

	config := &tcp.WaitConfig{}
	config.Address = tcpRegexp.ReplaceAllString(endpoint, "$1")
	config.Interval = time.Duration(interval * 1e6)
	config.Retries = retries
	config.Timeout = time.Duration(timeout * 1e6)

	config.OnAttempt = func(attempt int, config *tcp.WaitConfig, _ *error) {
		if attempt != 0 && !quiet {
			fmt.Fprintf(os.Stderr, "Waiting for %s (%d)...\n", config.Address, attempt)
		}
	}

	result, err := tcp.WaitTCPEndpoint(config)
	if err != nil {
		fail(quiet, "tcpwait error: %s", err)
	} else if !result.Success {
		fail(quiet, "could not reach \"%s\" after %fs", config.Address, result.Duration.Seconds())
	}

	succeed(quiet, "Reached \"%s\" in %fs", config.Address, result.Duration.Seconds())
}

func fail(quiet bool, format string, values ...interface{}) {
	if !quiet {
		fmt.Fprintf(os.Stderr, color.RedString("Error: "+format+"\n"), values...)
	}

	os.Exit(1)
}

func succeed(quiet bool, format string, values ...interface{}) {
	if !quiet {
		fmt.Fprintf(os.Stderr, color.GreenString(format+"\n", values...))
	}
}
