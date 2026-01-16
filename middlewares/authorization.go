package middlewares

import (
	"fmt"

	"github.com/PlanToPack/api-utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type AuthorizationPermissions struct {
	Admin interface{}
	Own   interface{}
}

func CheckAuthorization(neededPermissions AuthorizationPermissions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if rawPermissions, exists := ctx.Get(PermissionsTag); exists {
			permissions := (rawPermissions).([]interface{})
			if lo.Contains(permissions, neededPermissions.Admin) {
				ctx.Next()
				return
			} else if lo.Contains(permissions, neededPermissions.Own) {
				ctx.Set("userId", "fakeUserId") // TODO: Correct userId from the token claims (ctx.Keys)
				ctx.Next()
				return
			}
		}
		error := errors.NewForbiddenError(fmt.Errorf("Permission denied."))
		ctx.AbortWithStatusJSON(error.HTTPStatus, error)
		return
	}
}
