package retrier

import (
	"errors"
	"time"
)

type Periodic struct {
	maxAttempts int
	backoff     time.Duration
}

func NewPeriodicRetrier(attemps, waitTimeInMillis int) (*Periodic, error) {
	if attemps < 1 {
		return nil, errors.New("attempts should >= 1")
	}

	return &Periodic{
		maxAttempts: attemps,
		backoff:     time.Duration(waitTimeInMillis) * time.Millisecond,
	}, nil
}

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
