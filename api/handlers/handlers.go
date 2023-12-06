package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nixys/nxs-support-bot/ctx"
)

type RouteHandlerResponse struct {
	StatusCode int
	RAWData    any
	Message    string
}

type RouteHandlers struct {
	Handler       RouteHandler
	DataTransform RouteDataTransformHandler
}

type RouteHandler func(*ctx.Ctx, *gin.Context) RouteHandlerResponse
type RouteDataTransformHandler func(any, string) any

func RouteHandlerDefault(cc *ctx.Ctx, handler RouteHandlers) gin.HandlerFunc {
	return func(c *gin.Context) {

		if handler.Handler == nil {
			cc.Log.Warn("route handler not specified")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		r := handler.Handler(cc, c)

		var d interface{}
		if handler.DataTransform != nil {
			d = handler.DataTransform(r.RAWData, r.Message)
		} else {
			d = r.RAWData
		}

		if d != nil {
			c.JSON(r.StatusCode, d)
		} else {
			c.String(r.StatusCode, r.Message)
		}
	}
}
