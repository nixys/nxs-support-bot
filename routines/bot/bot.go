package bot

import (
	"context"

	appctx "github.com/nixys/nxs-go-appctx/v3"
	"github.com/nixys/nxs-support-bot/ctx"
	"github.com/sirupsen/logrus"
)

// Runtime executes the routine
func Runtime(app appctx.App) error {

	cc := app.ValueGet().(*ctx.Ctx)

	queueCh := make(chan error)
	updatesCh := make(chan error)

	queueCtx, queueCF := context.WithCancel(app.SelfCtx())
	defer queueCF()

	updatesCtx, updatesCF := context.WithCancel(app.SelfCtx())
	defer updatesCF()

	go cc.Bot.UpdatesGet(updatesCtx, updatesCh)
	go cc.Bot.Queue(queueCtx, queueCh)

	cc.Log.Debugf("telegram bot: successfully started")

	for {
		select {
		case <-app.SelfCtxDone():
			cc.Log.Debugf("telegram bot: shutdown")
			return nil
		case e := <-queueCh:
			if e != nil {
				cc.Log.WithFields(logrus.Fields{
					"details": e,
				}).Errorf("bot queue processing")
				app.Shutdown(e)
				return e
			}
		case e := <-updatesCh:
			if e != nil {
				cc.Log.WithFields(logrus.Fields{
					"details": e,
				}).Errorf("bot get updates")
				app.Shutdown(e)
				return e
			}
		}
	}
}
