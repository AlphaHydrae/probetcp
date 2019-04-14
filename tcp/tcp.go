// Package tcp provides the ProbeTCPEndpoint function to test if a TCP endpoint
// can be reached.
package tcp

import (
	"net"
	"time"
)

// ProbeTCPEndpoint tries to reach the configured TCP endpoint, returning either
// a ProbeResult if it worked, or an error if the endpoint could not be reached.
func ProbeTCPEndpoint(config *ProbeConfig) (*ProbeResult, error) {

	attempts := 0
	start := time.Now()

	var previousError *error
	var result *ProbeResult

	for config.Retries < 0 || attempts <= config.Retries {

		if attempts >= 1 && config.Interval >= 1 {
			time.Sleep(config.Interval)
		}

		if config.OnAttempt != nil {
			config.OnAttempt(attempts, config, previousError)
		}

		attempts++

		result = &ProbeResult{}

		conn, err := net.DialTimeout("tcp", config.Address, config.Timeout)
		if err != nil {
			if oerr, ok := err.(*net.OpError); ok {
				previousError = &err
				result.Error = oerr
			} else {
				return nil, err
			}
		} else {
			defer conn.Close()
			result.Connection = conn
			result.Error = nil
			result.Success = true
			break
		}
	}

	result.Attempts = attempts
	result.Duration = time.Now().Sub(start)

	return result, nil
}

// ProbeConfig is the configuration of a TCP probe. It is used to specify the
// TCP endpoint to call and can be used to set the retry behavior.
type ProbeConfig struct {

	// Address is the TCP address to call ("host:port", e.g. "localhost:5432",
	// "golang.org:http").
	Address string

	// Interval is the time to wait after a failed attempt before attempting a
	// retry. It is only used if Retries is set.
	Interval time.Duration

	// OnAttempt is a function that will be called on each TCP call during the
	// probe. It may be called several times if you specified Retries and a number
	// of calls fail.
	//
	// The function is called with the index of the attempt and the configuration
	// of the probe. The third err argument is nil if all goes well, but in the
	// case of a retry, it may contain the error that caused the previous call to
	// fail.
	OnAttempt func(attempt int, config *ProbeConfig, err *error)

	// Retries indicates how many times the TCP call should be retried if it
	// fails. If the first call(s) fail but one of the retries succeed, the probe
	// will be considered successful.
	Retries int

	// Timeout indicates the maximum time it should take to establish the TCP
	// connection. If it takes longer, the call will time out and the probe will
	// fail.
	Timeout time.Duration
}

// ProbeResult is the result of a TCP probe.  Its Success property indicates
// whether the probe was successful.
type ProbeResult struct {

	// Attempts indicates how many TCP calls were made. It will be 1 by default
	// but may be more if you specified Retries and some calls failed.
	Attempts int

	// Connection is the TCP connection that was established.
	Connection net.Conn

	// Duration indicates how much time it took to probe the endpoint.
	Duration time.Duration

	// Error will be nil if the probe succeeds, or it will contain the net.OpError
	// that caused the TCP call to fail.
	Error error

	// Success indicates whether the probe was successful.
	Success bool
}
