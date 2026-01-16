package middlewares

import (
	"github.com/PlanToPack/api-utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		detectedErrors := ctx.Errors.ByType(gin.ErrorTypeAny)

		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err
			log.Error().Err(err)
			parsedError, ok := err.(*errors.CustomError)
			if !ok {
				parsedError = errors.NewInternalServerError("UNKNOWN_ERROR", err)
			}
			ctx.AbortWithStatusJSON(parsedError.HTTPStatus, parsedError)
			return
		}
	}
}
