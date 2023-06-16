package bot

import (
	"context"

	"git.nixys.ru/apps/nxs-support-bot/ctx"
	appctx "github.com/nixys/nxs-go-appctx/v2"
)

// Runtime executes the routine
func Runtime(cr context.Context, appCtx *appctx.AppContext, crc chan interface{}) {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	queueCh := make(chan error)
	updatesCh := make(chan error)

	queueCtx, queueCF := context.WithCancel(cr)
	defer queueCF()

	updatesCtx, updatesCF := context.WithCancel(cr)
	defer updatesCF()

	go cc.Bot.UpdatesGet(updatesCtx, updatesCh)
	go cc.Bot.Queue(queueCtx, queueCh)

	for {
		select {
		case <-cr.Done():
			// Program termination.
			return
		case <-crc:
			// Updated context application data.
			// Set the new one in current goroutine.
			appCtx.Log().Errorf("context reload is not supported by runtime routine")
			appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
		case e := <-queueCh:
			if e != nil {
				appCtx.Log().Errorf("bot queue processing error: %v", e)
				appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
				return
			}
		case e := <-updatesCh:
			if e != nil {
				appCtx.Log().Errorf("bot get updates error: %v", e)
				appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
				return
			}
		}
	}
}
