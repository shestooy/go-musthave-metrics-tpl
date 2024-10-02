package utils

import (
	"errors"
	"net"
	"time"

	"github.com/avast/retry-go"
	"github.com/jackc/pgx/v5/pgconn"
)

func RetryDelay(n uint, _ error, _ *retry.Config) time.Duration {
	delays := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	if int(n) < len(delays) {
		return delays[n]
	}
	return delays[2]
}

func IsRetriableError(err error) bool {
	var connectErr *pgconn.ConnectError
	if errors.As(err, &connectErr) {
		return true
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}

	var netErr net.Error
	return errors.As(err, &netErr)

}
