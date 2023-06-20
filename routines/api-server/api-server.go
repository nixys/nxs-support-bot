package apiserver

import (
	"context"
	"net/http"
	"time"

	"github.com/nixys/nxs-support-bot/api"
	"github.com/nixys/nxs-support-bot/ctx"

	appctx "github.com/nixys/nxs-go-appctx/v2"
)

type httpServerContext struct {
	http.Server
	done chan interface{}
}

// Runtime executes the routine
func Runtime(cr context.Context, appCtx *appctx.AppContext, crc chan interface{}) {

	var err error

	s := servStart(appCtx)

	for {
		select {
		case <-cr.Done():
			// Program termination.
			err = servShutdown(s)
			if err != nil {
				appCtx.Log().Errorf("http server shutdown error: %v", err)
			}
			return
		case <-crc:
			// Updated context application data.
			// Set the new one in current goroutine.
			s, err = servRestart(appCtx, s)
			if err != nil {
				appCtx.Log().Errorf("http server reload error: %v", err)
				appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
				return
			}
		}
	}
}

func servStart(appCtx *appctx.AppContext) *httpServerContext {

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	s := &httpServerContext{
		Server: http.Server{
			Addr:         cc.Conf.API.Bind,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      api.RoutesSet(appCtx),
		},
		done: make(chan interface{}),
	}

	go func() {
		appCtx.Log().Debugf("server status: starting")
		if cc.Conf.API.TLS != nil {
			if err := s.ListenAndServeTLS(cc.Conf.API.TLS.CertFile, cc.Conf.API.TLS.KeyFie); err != nil {
				appCtx.Log().Debugf("server status: %v", err)
			}
		} else {
			if err := s.ListenAndServe(); err != nil {
				appCtx.Log().Debugf("server status: %v", err)
			}
		}
		s.done <- true
	}()

	return s
}

func servShutdown(s *httpServerContext) error {

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Shutdown the server
	if err := s.Shutdown(ctx); err != nil {
		return err
	}

	<-s.done
	return nil
}

func servRestart(appCtx *appctx.AppContext, s *httpServerContext) (*httpServerContext, error) {

	if err := servShutdown(s); err != nil {
		return nil, err
	}

	return servStart(appCtx), nil
}
