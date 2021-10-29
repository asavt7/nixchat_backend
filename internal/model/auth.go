package model

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// CachedTokens -  access tokens cached to auth Repo
type CachedTokens struct {
	AccessUID  string `json:"access"`
	RefreshUID string `json:"refresh"`
}

// Claims - JWT claims contains userinfo
type Claims struct {
	UserID uuid.UUID `json:"userId"`
	UID    uuid.UUID `json:"uid"`
	jwt.StandardClaims
}
