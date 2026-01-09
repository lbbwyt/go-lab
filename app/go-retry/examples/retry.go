package examples

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Error []error

func (e Error) Error() string {
	logWithNumber := make([]string, len(e))
	for i, l := range e {
		if l != nil {
			logWithNumber[i] = fmt.Sprintf("#%d: %s", i+1, l.Error())
		}
	}

	return fmt.Sprintf("All attempts fail:\n%s", strings.Join(logWithNumber, "\n"))
}

type RetryableFuncWithData[T any] func() (T, error)

func DoWithData[T any](ctx context.Context, retryableFunc RetryableFuncWithData[T], opts ...Option) (T, error) {
	var emptyT T
	var errs Error
	c := newDefaultConfig()
	for _, opt := range opts {
		opt(c)
	}
	backoff := c.delay

	for i := 0; i < c.attempts; i++ {
		t, err := retryableFunc()
		if err == nil {
			return t, nil
		}
		errs = append(errs, err)

		if i == c.attempts-1 {
			return emptyT, errs
		}
		sleep := c.delay
		//抖动
		jitter := time.Duration(rand.Int63n(int64(backoff)))
		sleep = backoff + jitter

		timer := time.NewTimer(sleep)
		select {
		case <-ctx.Done():
			timer.Stop()
			return emptyT, ctx.Err()
		case <-timer.C:
		}
		// 指数退避
		if c.backoff {
			backoff *= 2
			if backoff > c.maxDelay {
				backoff = c.maxDelay
			}
		}

	}
	return emptyT, errs
}

func Do(ctx context.Context, fn func() error, opts ...Option) error {
	c := newDefaultConfig()
	for _, opt := range opts {
		opt(c)
	}
	retryableFuncWithData := func() (any, error) {
		return nil, fn()
	}
	_, err := DoWithData(ctx, retryableFuncWithData, opts...)
	return err
}
