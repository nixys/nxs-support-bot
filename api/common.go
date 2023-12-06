package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.WithFields(logrus.Fields{
			"type":      "accesslog",
			"remote":    c.RemoteIP(),
			"method":    c.Request.Method,
			"url":       c.Request.RequestURI,
			"code":      c.Writer.Status(),
			"userAgent": c.Request.UserAgent(),
		}).Info("request processed")
	}
}

func RequestSizeLimiter(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > limit {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//c.Writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, X-Auth-Health-Key, X-Auth-Key, If-Modified-Since, Cache-Control, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
	}
}
