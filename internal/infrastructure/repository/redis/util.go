package rredis

import (
	"strings"
)

func buildKey(args ...string) string {
	return strings.Join(args, "_")
}
