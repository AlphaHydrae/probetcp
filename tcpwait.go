// The tcpwait command checks whether a TCP endpoint can be reached.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"syscall"
	"time"

	"github.com/alphahydrae/tcpwait/tcp"
	"github.com/fatih/color"
	flag "github.com/spf13/pflag"
)

const usageHeader = `%s waits for TCP endpoints to be reachable.

Usage:
  %s [OPTION...] ENDPOINT... [--] [EXEC...]

Options:
`

const usageFooter = `
Examples:
  Wait for a website:
    tcpwait google.com:80
  Wait for a MySQL database (10 attempts every 2 seconds):
    tcpwait -r 9 -i 2000 tcp://localhost:3306
  Wait for multiple endpoints:
    tcpwait github.com:22 github.com:80 github.com:443
  Execute a command after an endpoint is reached:
    tcpwait db.example.com:5432 -- rails server
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

	if interval < 0 {
		fail(1, quiet, "the \"interval\" option must be greater than or equal to zero")
	} else if retries < 0 {
		fail(1, quiet, "the \"retries\" option must be greater than or equal to zero")
	} else if timeout <= 0 {
		fail(1, quiet, "the \"timeout\" option must be greater than zero")
	}

	terminator := -1
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "--" {
			terminator = i
			break
		}
	}

	var endpoints []string
	var execCommand string
	var execArgs []string
	if terminator >= 0 && terminator < len(os.Args)-1 {
		endpoints = flag.Args()[0 : len(flag.Args())-(len(os.Args)-terminator-1)]

		var err error
		execCommand, err = exec.LookPath(os.Args[terminator+1])
		if err != nil {
			fail(10, quiet, "could not find command \"%s\"", os.Args[terminator+1])
		}

		execArgs = append([]string{execCommand}, os.Args[terminator+2:len(os.Args)]...)
	} else {
		endpoints = flag.Args()
	}

	if len(endpoints) == 0 || endpoints[0] == "" {
		fail(1, quiet, "an endpoint to wait for must be given as an argument (e.g. \"tcp://localhost:3306\")")
	}

	ch := make(chan *waitResult)
	tcpRegexp := regexp.MustCompile(tcpTargetRegexp)

	for i := 0; i < len(endpoints); i++ {
		endpoint := endpoints[i]

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

		go wait(config, ch)
	}

	for i := 0; i < len(endpoints); i++ {
		result := <-ch
		if result.error != nil {
			fail(2, quiet, "tcpwait error: %s", result.error)
		} else if !result.result.Success {
			fail(3, quiet, "could not reach \"%s\" after %fs", result.config.Address, result.result.Duration.Seconds())
		} else {
			succeed(quiet, "Reached \"%s\" in %fs", result.config.Address, result.result.Duration.Seconds())
		}
	}

	if execCommand != "" {
		err := syscall.Exec(execCommand, execArgs, os.Environ())
		if err != nil {
			fail(11, quiet, "could not execute command \"%s\" with arguments %s", execCommand, execArgs)
		}
	}
}

func wait(config *tcp.WaitConfig, ch chan *waitResult) {
	result, err := tcp.WaitTCPEndpoint(config)

	chResult := &waitResult{}
	chResult.config = config
	chResult.result = result
	chResult.error = err

	ch <- chResult
}

func fail(code int, quiet bool, format string, values ...interface{}) {
	if !quiet {
		fmt.Fprintf(os.Stderr, color.RedString("Error: "+format+"\n"), values...)
	}

	os.Exit(code)
}

func succeed(quiet bool, format string, values ...interface{}) {
	if !quiet {
		fmt.Fprintf(os.Stderr, color.GreenString(format+"\n", values...))
	}
}

type waitResult struct {
	config *tcp.WaitConfig
	result *tcp.WaitResult
	error  error
}
