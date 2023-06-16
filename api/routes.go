package api

import (
	"git.nixys.ru/apps/nxs-support-bot/api/endpoints"
	"git.nixys.ru/apps/nxs-support-bot/ctx"

	"github.com/gin-gonic/gin"
	appctx "github.com/nixys/nxs-go-appctx/v2"
)

func RoutesSet(appCtx *appctx.AppContext) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(endpoints.Logger(appCtx))
	router.Use(endpoints.CORSMiddleware())

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	v1 := router.Group("/v1")
	{
		v1.Use(endpoints.RequestSizeLimiter(appCtx))

		redmine := v1.Group("/redmine")
		{
			redmine.Use(endpoints.AuthorizeRedmine(cc.Conf.API.RedmineSecretToken))

			redmine.POST("", endpoints.RouteHandlerDefault(appCtx, endpoints.RouteHandlers{
				Handler: endpoints.Redmine,
			}))
		}
	}

	return router
}
