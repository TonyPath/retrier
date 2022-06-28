package retrier

import (
	"errors"
)

var ErrLimitExceeded = errors.New("retrier have reached the maximum attempts")

type WrappedTask func() (interface{}, error)

type IsRetryAble func(err error) bool

type Retrier interface {
	Retry(WrappedTask, IsRetryAble) (interface{}, error)
}
