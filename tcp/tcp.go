// Package tcp provides the WaitTCPEndpoint function to wait until a TCP
// endpoint can be reached.
package tcp

import (
	"net"
	"time"
)

// WaitTCPEndpoint tries to reach the configured TCP endpoint, returning either
// a WaitResult if it worked, or an error if the endpoint could not be reached.
func WaitTCPEndpoint(config *WaitConfig) (*WaitResult, error) {

	attempts := 0
	start := time.Now()

	var previousError *error
	var result *WaitResult

	for config.Retries < 0 || attempts <= config.Retries {

		if attempts >= 1 && config.Interval >= 1 {
			time.Sleep(config.Interval)
		}

		if config.OnAttempt != nil {
			config.OnAttempt(attempts, config, previousError)
		}

		attempts++

		result = &WaitResult{}

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

// WaitConfig is the configuration to wait for a TCP endpoint to be reachable.
// It is used to specify the address to call and retry behavior.
type WaitConfig struct {

	// Address is the TCP address to try to reach ("host:port", e.g.
	// "localhost:5432", "golang.org:http").
	Address string

	// Interval is the time to wait after a failed attempt before retrying. It is
	// only used if Retries is set.
	Interval time.Duration

	// OnAttempt is a function that will be called each time an attempt is made to
	// reach the TCP endpoint.  It may be called several times if you specified
	// Retries and a number of calls fail.
	//
	// The function is called with the index of the attempt and the waiting
	// configuration. The third err argument is nil if all goes well, but in the
	// case of a retry, it may contain the error that caused the previous call to
	// fail.
	OnAttempt func(attempt int, config *WaitConfig, err *error)

	// Retries indicates how many times to attempt to reach the TCP endpoint again
	// if the first attempt fails. The wait is considered successful as soon as
	// one attempt succeeds, regardless of previous failures.
	Retries int

	// Timeout indicates the maximum time to wait to establish the TCP connection
	// on each attempt.  If it takes longer, the call will time out and the
	// attempt will fail.
	Timeout time.Duration
}

// WaitResult is the result of waiting to reach a TCP endpoint. Its Success
// property indicates whether the endpoint was reached.
type WaitResult struct {

	// Attempts indicates how many attempts were made to reach the endpoint.  It
	// will be 1 by default but may be more if you specified Retries and some
	// calls failed.
	Attempts int

	// Connection is the TCP connection that was established.
	Connection net.Conn

	// Duration indicates how much time it took to reach the endpoint.
	Duration time.Duration

	// Error will be nil if the endpoint has been reached, or it will contain the
	// net.OpError that caused the TCP call to fail.
	Error error

	// Success indicates whether the endpoint could be reached (the first time or
	// on a subsequent retry attempt).
	Success bool
}
