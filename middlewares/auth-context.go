package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ginerator/base/errors"
	user "github.com/ginerator/base/model/users"
)

func buildSystemUserFromAuth(c ClientMetadata) user.User {
	var source string = c.ClientName

	return user.User{
		Source: &source,
		Type:   user.UserTypeSystem,
	}
}

func buildUserFromAuth(u UserMetadata) user.User {
	var email string = u.Email
	var id string = u.UserId

	return user.User{
		Email: &email,
		Id:    &id,
		Type:  user.UserTypePerson,
	}
}

func SetUserContext(ctx *gin.Context) {
	claims, exists := ctx.Get(CustomClaimsTag)
	if !exists {
		error := errors.NewUnauthorizedError(fmt.Errorf("User claims are invalid"))
		ctx.AbortWithStatusJSON(error.HTTPStatus, error)
	}
	customClaims := claims.(CustomClaims)
	clientType := customClaims.ClientType

	switch clientType {
	case ClientTypeUser:
		ctx.Set(user.ContextTagUser, buildUserFromAuth(customClaims.UserMetadata))
	default:
		ctx.Set(user.ContextTagUser, buildSystemUserFromAuth(customClaims.ClientMetadata))
	}
}
