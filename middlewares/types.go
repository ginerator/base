package middlewares

import (
	"github.com/golang-jwt/jwt/v4"
)

type ClientType string

const (
	ClientTypeSystem ClientType = "MACHINE"
	ClientTypeUser              = "USER"
)

type ClientMetadata struct {
	ClientName string `json:"client_name"`
}

type UserMetadata struct {
	UserId string `json:"account_id"`
	Email  string `json:"email"`
}

type CustomClaims struct {
	ClientMetadata ClientMetadata `json:"client_metadata,omitempty"`
	ClientType     string         `json:"client_type"`
	Roles          []string       `json:"roles"`
	UserMetadata   UserMetadata   `json:"user_metadata,omitempty"`
}

type Auth0TokenClaims struct {
	CustomClaims    CustomClaims  `json:"https://hear.com"`
	Permissions     []interface{} `json:"permissions"`
	Scope           string        `json:"scope"`
	AuthorizedParty string        `json:"azp"`
	jwt.RegisteredClaims
}
