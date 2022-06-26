package utils

import (
	"context"
	"sync"
	"time"
)

func WaitGroupTimeOut(timeout time.Duration) (*sync.WaitGroup, context.Context, context.CancelFunc) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	go func() {
		wg.Wait()
		cancel()
	}()
	return &wg, ctx, cancel
}
