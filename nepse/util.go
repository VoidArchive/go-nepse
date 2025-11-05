package nepse

import "time"

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
