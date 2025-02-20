package try

import (
	"errors"
	"fmt"
	"time"
)

// Number of max retry allowed
var maxRetries = 10

// ErrMaxRetriesReached is the error returned when maxRetries exceeded.
var ErrMaxRetriesReached = errors.New("exceeded retry limit")

// retryableFn is the function signature in a retry.
type retryableFn func(attempt int) (retry bool, err error)

// Do runs fn and retries if fn fails. Do exits when either error returned by fn is nil,
// retry returned by fn is false, or maxRetries exceeded.
func Do(fn retryableFn) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > maxRetries {
			return ErrMaxRetriesReached
		}
	}
	return err
}

func DoBackOff(fn retryableFn, dur time.Duration) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > maxRetries {
			return fmt.Errorf("%s %s", ErrMaxRetriesReached, err)
		}
		time.Sleep(dur)
	}
	return err
}
