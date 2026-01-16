package validators

import (
	"github.com/PlanToPack/api-utils/utils"
	"github.com/gin-gonic/gin"
	lo "github.com/samber/lo"
)

func IsValidQuery(ctx *gin.Context, s interface{}) (bool, []string) {
	queryParams := ctx.Request.URL.Query()
	allowedParams := utils.GetStructKeys(s)

	queryKeys := lo.Keys[string, []string](queryParams)
	unknownKeys, _ := lo.Difference(queryKeys, allowedParams)

	isValid := true
	if len(unknownKeys) > 0 {
		isValid = false
	}

	return isValid, unknownKeys
}
