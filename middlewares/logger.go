package middlewares

import (
	ginLogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return ginLogger.SetLogger(ginLogger.WithSkipPath([]string{"/sys/health"}))
}
