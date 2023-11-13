package appcache

import (
	"time"

	appctx "github.com/nixys/nxs-go-appctx/v3"
	"github.com/nixys/nxs-support-bot/ctx"
	"github.com/sirupsen/logrus"
)

func Runtime(app appctx.App) error {

	cc := app.ValueGet().(*ctx.Ctx)

	if err := cc.Cache.C.Update(); err != nil {
		cc.Log.WithFields(logrus.Fields{
			"details": err,
		}).Errorf("cache runtime")
		return err
	}
	cc.Log.Debugf("cache: successfully updated")

	timer := time.NewTimer(cc.Cache.TTL)

	for {
		select {
		case <-app.SelfCtxDone():
			cc.Log.Debugf("cache: shutdown")
			return nil
		case <-timer.C:
			if err := cc.Cache.C.Update(); err != nil {
				cc.Log.WithFields(logrus.Fields{
					"details": err,
				}).Errorf("cache runtime")
				return err
			}
			cc.Log.Debugf("cache: successfully updated")

			timer.Reset(cc.Cache.TTL)
		}
	}
}
