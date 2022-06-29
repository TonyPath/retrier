/*
Package retrier provides support for retry the execution of a task

Example Usage

	var ErrRetryAble = errors.New("retryable error")

	retrier, err := NewPeriodicRetrier(10, 1000)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	resp, err := retrier.Retry(func() (interface{}, error) {
			resp, err := http.Get("URL")
			return resp, err
		}, func(err error) bool {
			if errors.Is(err, ErrRetryAble) {
				return true
			}
			return false
		})
*/
package retrier
