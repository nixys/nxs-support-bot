package appcache

import (
	"context"
	"time"

	appctx "github.com/nixys/nxs-go-appctx/v2"

	"github.com/nixys/nxs-support-bot/ctx"
)

// Runtime executes the routine
func Runtime(c context.Context, appCtx *appctx.AppContext, crc chan interface{}) {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	timer := time.NewTimer(cc.Cache.TTL)

	if err := cacheUpdate(appCtx); err != nil {
		appCtx.Log().Errorf("cache runtime error: %v", err)
	}

	for {
		select {
		case <-timer.C:
			// uwatch iterate
			if err := cacheUpdate(appCtx); err != nil {
				appCtx.Log().Errorf("cache runtime error: %v", err)
			}
			timer.Reset(cc.Cache.TTL)
		case <-c.Done():
			// Program termination.
			// Write "Done" to log and complete the current goroutine.
			appCtx.Log().Info("cache done")
			return
		case <-crc:
			// Updated context application data.
			// Set the new one in current goroutine.
			appCtx.Log().Info("cache routine reload")
		}
	}
}

func cacheUpdate(appCtx *appctx.AppContext) error {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	if err := cc.Cache.C.Update(); err != nil {
		appCtx.Log().Errorf("cache runtime error: %v", err)
		return err
	}

	return nil
}
