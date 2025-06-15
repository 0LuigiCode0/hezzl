package rredis

import (
	"crypto/sha1"
	"strings"
)

func hashKey(args ...string) string {
	hash := sha1.Sum([]byte(strings.Join(args, "_")))
	return string(hash[:])
}
