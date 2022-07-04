package retrier

import (
	"errors"
)

// ErrLimitExceeded indicates that the retrier has reached the limit of attempts.
var ErrLimitExceeded = errors.New("retrier has reached the maximum attempts")

// WrappedTask is function designed to be used as a wrapper of the task to be executed.
type WrappedTask func() (any, error)

// IsRetryAble is a function designed to decide if the execution of a task should be retried or not.
type IsRetryAble func(err error) bool

// Retrier defines an interface for providing the implementation details for
// (re)executing a given task.
type Retrier interface {
	Retry(WrappedTask, IsRetryAble) (any, error)
}
