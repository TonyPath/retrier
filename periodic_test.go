package retrier

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPeriodicRetrier_Retry(t *testing.T) {

	var (
		errNonRetryAble = errors.New("non-retryable error")
		errRetryAble    = errors.New("retryable error")

		calls       = 0
		maxAttempts = 3
		backoff     = 2 * time.Millisecond
	)

	tests := map[string]struct {
		task        func() (res interface{}, err error)
		wantCall    int
		wantResult  interface{}
		wantErr     bool
		wantErrType error
	}{
		"succeeded immediately": {
			task: func() (res interface{}, err error) {
				calls++
				return "task response", nil
			},
			wantCall:    1,
			wantResult:  "task response",
			wantErr:     false,
			wantErrType: nil,
		},
		"succeeded at 3rd attempt": {
			task: func() (res interface{}, err error) {
				calls++
				if calls == 3 {
					return []string{"task response"}, nil
				}
				return nil, errRetryAble
			},
			wantCall:    3,
			wantResult:  []string{"task response"},
			wantErr:     false,
			wantErrType: nil,
		},
		"reach max attempts": {
			task: func() (res interface{}, err error) {
				calls++
				return "", errRetryAble
			},
			wantCall:    3,
			wantErr:     true,
			wantErrType: ErrLimitExceeded,
		},
		"non-retryable- stop retrying": {
			task: func() (res interface{}, err error) {
				calls++
				return nil, errNonRetryAble
			},
			wantCall:    1,
			wantErr:     true,
			wantErrType: errNonRetryAble,
		},
	}

	for name, tcase := range tests {
		t.Run(name, func(t *testing.T) {
			calls = 0

			retrier := Periodic{
				maxAttempts: maxAttempts,
				backoff:     backoff,
			}

			res, err := retrier.Retry(func() (interface{}, error) {
				res, err := tcase.task()
				return res, err
			}, func(err error) bool {
				if errors.Is(err, errRetryAble) {
					return true
				}
				return false
			})

			if tcase.wantErr {
				require.ErrorIs(t, err, tcase.wantErrType)
			} else {
				require.Nil(t, err)
			}
			require.Equal(t, tcase.wantResult, res)
			require.Equal(t, tcase.wantCall, calls)
		})
	}
}

func TestNewPeriodicRetrier(t *testing.T) {
	type args struct {
		attemps          int
		waitTimeInMillis int
	}
	tests := []struct {
		name    string
		args    args
		want    *Periodic
		wantErr bool
	}{
		{
			name: "",
			args: args{
				attemps:          3,
				waitTimeInMillis: 1000,
			},
			want: &Periodic{
				maxAttempts: 3,
				backoff:     1000 * time.Millisecond,
			},
			wantErr: false,
		},
		{
			name: "a",
			args: args{
				attemps:          0,
				waitTimeInMillis: 1000,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPeriodicRetrier(tt.args.attemps, tt.args.waitTimeInMillis)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPeriodicRetrier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPeriodicRetrier() got = %v, want %v", got, tt.want)
			}
		})
	}
}
