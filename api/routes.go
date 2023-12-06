package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nixys/nxs-support-bot/api/endpoints"
	endpointsv2 "github.com/nixys/nxs-support-bot/api/endpoints/v2"
	"github.com/nixys/nxs-support-bot/api/handlers"
	"github.com/nixys/nxs-support-bot/ctx"
)

func RoutesSet(cc *ctx.Ctx) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(Logger(cc.Log))
	router.Use(CORSMiddleware())
	router.Use(RequestSizeLimiter(cc.API.ClientMaxBodySizeBytes))

	v1 := router.Group("/v1")
	{
		redmine := v1.Group("/redmine")
		{
			redmine.Use(endpoints.AuthorizeRedmine(cc.API.RedmineSecretToken))

			redmine.POST("", handlers.RouteHandlerDefault(cc, handlers.RouteHandlers{
				Handler: endpoints.Redmine,
			}))
		}
	}

	v2 := router.Group("/v2")
	{
		redmine := v2.Group("/redmine")
		{
			redmine.Use(endpointsv2.AuthorizeRedmine(cc.API.RedmineSecretToken))

			redmine.POST("/created", handlers.RouteHandlerDefault(cc, handlers.RouteHandlers{
				Handler: endpointsv2.RedmineCreated,
			}))

			redmine.POST("/updated", handlers.RouteHandlerDefault(cc, handlers.RouteHandlers{
				Handler: endpointsv2.RedmineUpdated,
			}))
		}
	}

	return router
}
