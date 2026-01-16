package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ginerator/base/utils"
)

func AttachSysRoutes(router *gin.Engine, appName string, appStateManager *utils.AppStateManager) *gin.RouterGroup {
	sys := router.Group("/sys")
	sys.GET("/health", func(ctx *gin.Context) {
		_, err := appStateManager.DependenciesConnected()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"name":   appName,
				"status": "DOWN",
				"error":  err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"name":   appName,
			"status": "UP",
		})
	})
	return sys
}
