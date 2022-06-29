package retrier

import (
	"errors"
	"time"
)

// Periodic represents an implementation of Retrier interface.
type Periodic struct {
	maxAttempts int
	backoff     time.Duration
}

// NewPeriodicRetrier returns a new Periodic value for use.
func NewPeriodicRetrier(attemps, waitTimeInMillis int) (*Periodic, error) {
	if attemps < 1 {
		return nil, errors.New("attempts should >= 1")
	}

	return &Periodic{
		maxAttempts: attemps,
		backoff:     time.Duration(waitTimeInMillis) * time.Millisecond,
	}, nil
}

// Retry performs the execution of a task and applies a retry policy.
func (r *Periodic) Retry(wrappedTask WrappedTask, isRetryAble IsRetryAble) (interface{}, error) {
	var attempts = 0

	for {
		attempts++

		res, err := wrappedTask()
		if err == nil {
			return res, nil
		}

		if isRetryAble(err) {
			if attempts == r.maxAttempts {
				return nil, ErrLimitExceeded
			}
			time.Sleep(r.backoff)
		} else {
			return nil, err
		}
	}
}
