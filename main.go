package main

import (
	"os"
	"syscall"

	appctx "github.com/nixys/nxs-go-appctx/v3"

	"github.com/nixys/nxs-support-bot/ctx"
	"github.com/nixys/nxs-support-bot/misc"
	apiserver "github.com/nixys/nxs-support-bot/routines/api-server"
	appcache "github.com/nixys/nxs-support-bot/routines/app-cache"
	"github.com/nixys/nxs-support-bot/routines/bot"
)

func main() {

	err := appctx.Init(nil).
		RoutinesSet(
			map[string]appctx.RoutineParam{
				"cache": {
					Handler: appcache.Runtime,
				},
				"apiserver": {
					Handler: apiserver.Runtime,
				},
				"bot": {
					Handler: bot.Runtime,
				},
			},
		).
		ValueInitHandlerSet(ctx.AppCtxInit).
		SignalsSet([]appctx.SignalsParam{
			{
				Signals: []os.Signal{
					syscall.SIGTERM,
				},
				Handler: sigHandlerTerm,
			},
		}).
		Run()
	if err != nil {
		switch err {
		case misc.ErrArgSuccessExit:
			os.Exit(0)
		default:
			os.Exit(1)
		}
	}
}

func sigHandlerTerm(sig appctx.Signal) {
	sig.Shutdown(nil)
}
