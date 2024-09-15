package utils

import (
	"github.com/avast/retry-go"
	"time"
)

func RetryDelay(n uint, _ error, _ *retry.Config) time.Duration {
	delays := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	if int(n) < len(delays) {
		return delays[n]
	}
	return delays[2]
}
