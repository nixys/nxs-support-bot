package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	appctx "github.com/nixys/nxs-go-appctx/v2"
	"github.com/sirupsen/logrus"

	"git.nixys.ru/apps/nxs-support-bot/ctx"
	apiserver "git.nixys.ru/apps/nxs-support-bot/routines/api-server"
	appcache "git.nixys.ru/apps/nxs-support-bot/routines/app-cache"
	"git.nixys.ru/apps/nxs-support-bot/routines/bot"
)

func main() {

	// Read command line arguments
	args := ctx.ArgsRead()

	appCtx, err := appctx.ContextInit(appctx.Settings{
		CustomContext:    &ctx.Ctx{},
		Args:             &args,
		CfgPath:          args.ConfigPath,
		TermSignals:      []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		ReloadSignals:    []os.Signal{syscall.SIGHUP},
		LogrotateSignals: []os.Signal{syscall.SIGUSR1},
		LogFormatter:     &logrus.JSONFormatter{},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	appCtx.Log().Info("program started")

	// main() body function
	defer appCtx.MainBodyGeneric()

	// Create main context
	c := context.Background()

	// Create app cache routine
	appCtx.RoutineCreate(c, appcache.Runtime)

	// Create bot routine
	appCtx.RoutineCreate(c, bot.Runtime)

	// Create API server routine
	appCtx.RoutineCreate(c, apiserver.Runtime)
}
