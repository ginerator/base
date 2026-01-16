package middlewares

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/PlanToPack/api-utils/errors"
	user "github.com/PlanToPack/api-utils/model/users"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

const (
	CustomClaimsTag = "customClaims"
	PermissionsTag  = "permissions"
	Token           = "token"
)

func NewJWKSProvider(jwksBaseURL string) *keyfunc.JWKS {
	jwksURL, err := url.Parse(fmt.Sprintf("%s/.well-known/jwks.json", jwksBaseURL))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse Auth0 jwks URL.")
	}

	options := keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshErrorHandler: func(err error) {
			log.Fatal().Err(err).Msg("Failed to refresh JWKS.")
		},
	}

	jwks, err := keyfunc.Get(jwksURL.String(), options)
	if err != nil {
		log.Fatal().Err(err).Str("url", jwksURL.String()).Msg("Failed to create JWKS from resource at the given URL.")
	}
	return jwks
}

func Authenticate(provider *keyfunc.JWKS) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.Request.Header.Get("Authorization")

		jwtToken, err := stripBearerToken(authorizationHeader)
		if err != nil {
			error := errors.NewUnauthorizedError(err)
			ctx.AbortWithStatusJSON(error.HTTPStatus, error)
			return
		}

		claims := new(Auth0TokenClaims)

		token, err := jwt.ParseWithClaims(jwtToken, claims, provider.Keyfunc)
		if err != nil || !token.Valid {
			error := errors.NewUnauthorizedError(fmt.Errorf("Error parsing token: %v.", err))
			ctx.AbortWithStatusJSON(error.HTTPStatus, error)
			return
		}

		ctx.Set(PermissionsTag, claims.Permissions)
		ctx.Set(CustomClaimsTag, claims.CustomClaims)
		ctx.Set(user.ContextTagToken, authorizationHeader)
		SetUserContext(ctx)
		ctx.Next()
	}
}

func stripBearerToken(token string) (string, error) {
	if len(token) > 6 && strings.ToUpper(token[0:7]) == "BEARER " {
		return token[7:], nil
	}
	return "", fmt.Errorf("Invalid token.")
}
