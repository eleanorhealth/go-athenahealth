package ratelimiter

import "errors"

var ErrRateExceeded = errors.New("rate limit exceeded")
