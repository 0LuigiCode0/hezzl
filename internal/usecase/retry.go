package usecase

import (
	"context"
	"time"
)

func (u *_usecase) withRetry(count int, f func(ctx context.Context) error) {
	u.wg.Add(1)
	go func() {
		defer u.wg.Done()
		for i := range count {
			select {
			case <-u.ctx.Done():
				return
			case <-time.After(time.Second * time.Duration(i*i)):
				err := f(u.ctx)
				if err == nil {
					return
				}
			}
		}
	}()
}
