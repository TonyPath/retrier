# Go Retrier

[![codecov](https://codecov.io/gh/TonyPath/retrier/branch/master/graph/badge.svg?token=MNSBQIUJBK)](https://codecov.io/gh/TonyPath/retrier)
[![Go Report Card](https://goreportcard.com/badge/github.com/TonyPath/retrier)](https://goreportcard.com/report/github.com/TonyPath/retrier)
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/TonyPath/retrier)](https://github.com/TonyPath/retrier)
[![GoDoc](https://godoc.org/github.com/TonyPath/retrier?status.svg)](https://godoc.org/github.com/TonyPath/retrier)

Simple retrier in Go

## Examples

```go
maxAttempts := 5
backoffInMillis := 1000
retrier, _ := NewPeriodicRetrier(maxAttempts, backoffInMillis)

res, err := retrier.Retry(
	func() (interface{}, error) {
            res, err := http.Get(...)
            return res, err
	}, 
	func(err error) bool {
            var errDNS *net.DNSError
            if errors.As(err, &errDNS) && errDNS.IsTemporary {
                return true
            }
            return false
        }
)
```