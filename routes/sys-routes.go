package routes

import (
	"net/http"

	"github.com/PlanToPack/api-utils/utils"
	"github.com/gin-gonic/gin"
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
