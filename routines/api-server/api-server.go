package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	appctx "github.com/nixys/nxs-go-appctx/v3"
	"github.com/nixys/nxs-support-bot/api"
	"github.com/nixys/nxs-support-bot/ctx"
	"github.com/sirupsen/logrus"
)

type httpServerContext struct {
	http.Server
	done chan interface{}
}

func Runtime(app appctx.App) error {

	cc := app.ValueGet().(*ctx.Ctx)

	s := servStart(cc)

	for {
		select {
		case <-app.SelfCtxDone():

			c, f := context.WithTimeout(app.SelfCtx(), 1*time.Second)
			defer f()

			cc.Log.Debugf("api: shutting down")

			err := servShutdown(c, s)
			if err != nil {

				cc.Log.WithFields(logrus.Fields{
					"details": err,
				}).Errorf("api: shutdown")

				err = fmt.Errorf("api: %w", err)
			}
			return err
		}
	}
}

func servStart(cc *ctx.Ctx) *httpServerContext {

	s := &httpServerContext{
		Server: http.Server{
			Addr:         cc.API.Bind,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      api.RoutesSet(cc),
		},
		done: make(chan interface{}),
	}

	go func() {
		cc.Log.Debugf("api: starting")
		if cc.API.TLS != nil {
			if err := s.ListenAndServeTLS(cc.API.TLS.CertFile, cc.API.TLS.KeyFie); err != nil {
				cc.Log.WithFields(logrus.Fields{
					"details": err,
				}).Debugf("api: server listen tls")
			}
		} else {
			if err := s.ListenAndServe(); err != nil {
				cc.Log.WithFields(logrus.Fields{
					"details": err,
				}).Debugf("api: server listen")
			}
		}
		s.done <- true
	}()

	return s
}

func servShutdown(c context.Context, s *httpServerContext) error {

	//Shutdown the server
	if err := s.Shutdown(c); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	<-s.done
	return nil
}
