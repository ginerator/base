package utils

import "github.com/gin-gonic/gin"

func GetUserId(ctx *gin.Context) *string {
	if userId, exists := ctx.Keys["userId"]; exists {
		if strUserId, ok := userId.(string); ok {
			return &strUserId
		}
	}
	return nil
}
