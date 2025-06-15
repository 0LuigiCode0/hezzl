package ihttp

import (
	"context"
	"time"
)

func (h *_handler) withRetry(count int, f func(ctx context.Context) error) {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		for i := range count {
			select {
			case <-h.ctx.Done():
				return
			case <-time.After(time.Second * time.Duration(i*i)):
				err := f(h.ctx)
				if err == nil {
					return
				}
			}
		}
	}()
}
