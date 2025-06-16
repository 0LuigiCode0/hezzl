package utils

import "time"

var IsTimeMock bool

func TimeNow() time.Time {
	if IsTimeMock {
		return time.Time{}
	}
	return time.Now()
}
