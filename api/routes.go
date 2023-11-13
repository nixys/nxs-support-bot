package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nixys/nxs-support-bot/api/endpoints"
	"github.com/nixys/nxs-support-bot/ctx"
)

func RoutesSet(cc *ctx.Ctx) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(endpoints.Logger(cc.Log))
	router.Use(endpoints.CORSMiddleware())

	v1 := router.Group("/v1")
	{
		v1.Use(endpoints.RequestSizeLimiter(cc.API.ClientMaxBodySizeBytes))

		redmine := v1.Group("/redmine")
		{
			redmine.Use(endpoints.AuthorizeRedmine(cc.Conf.API.RedmineSecretToken))

			redmine.POST("", endpoints.RouteHandlerDefault(cc, endpoints.RouteHandlers{
				Handler: endpoints.Redmine,
			}))
		}
	}

	return router
}
