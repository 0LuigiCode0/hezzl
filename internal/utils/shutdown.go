package utils

import (
	"slices"
	"sync"
)

var shutdownQuery []func()

func AddShutdown(f func()) {
	shutdownQuery = append(shutdownQuery, f)
}

var shutdownOnce = sync.OnceFunc(func() {
	for _, f := range slices.Backward(shutdownQuery) {
		f()
	}
})

func Shutdown() { shutdownOnce() }
