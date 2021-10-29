package auth

import (
	"errors"
	"github.com/asavt7/nixchat_backend/internal/config"
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/asavt7/nixchat_backend/internal/repos"
	"github.com/asavt7/nixchat_backend/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type AuthorizationServiceImpl struct {
	cfg        *config.AuthConfig
	userRepo   repos.UserRepo
	tokenStore repos.TokenKeeper
}

func NewAuthorizationServiceImpl(cfg *config.AuthConfig, userRepo repos.UserRepo, tokenStore repos.TokenKeeper) *AuthorizationServiceImpl {
	return &AuthorizationServiceImpl{cfg: cfg, userRepo: userRepo, tokenStore: tokenStore}
}

func (s *AuthorizationServiceImpl) GetAccessSigningKey() string {
	return s.cfg.AccessSecret
}

func (s *AuthorizationServiceImpl) Logout(accessTokenClaims *model.Claims) error {
	userID := accessTokenClaims.UserID.String()
	return s.tokenStore.Delete(userID)
}

func (s *AuthorizationServiceImpl) ValidateRefreshToken(accessTokenClaims *model.Claims) error {
	return s.validateToken(accessTokenClaims, func(tokens model.CachedTokens) string {
		return tokens.RefreshUID
	})
}

func (s *AuthorizationServiceImpl) ValidateAccessToken(accessTokenClaims *model.Claims) error {
	return s.validateToken(accessTokenClaims, func(tokens model.CachedTokens) string {
		return tokens.AccessUID
	})
}

func (s *AuthorizationServiceImpl) validateToken(accessTokenClaims *model.Claims, tokenFromCashedFunc func(model.CachedTokens) string) error {
	userID := accessTokenClaims.UserID.String()
	clientTokenUID := accessTokenClaims.UID

	cached, err := s.tokenStore.Get(userID)
	if err != nil {
		return err
	}
	if tokenFromCashedFunc(cached) != clientTokenUID.String() {
		return errors.New("invalid token")
	}
	return nil
}

// CheckUserCredentials
func (s *AuthorizationServiceImpl) CheckUserCredentials(username string, password string) (model.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return user, err
	}
	err = utils.CheckPassword(user.PasswordHash, password)
	return user, err
}

func (s *AuthorizationServiceImpl) IsNeedToRefresh(claims *model.Claims) bool {
	return time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15*time.Minute
}

func (s *AuthorizationServiceImpl) ParseAccessTokenToClaims(token string) (*model.Claims, error) {
	return s.parseTokenToClaims(token, []byte(s.cfg.RefreshSecret))
}

func (s *AuthorizationServiceImpl) ParseRefreshTokenToClaims(token string) (*model.Claims, error) {
	return s.parseTokenToClaims(token, []byte(s.cfg.RefreshSecret))
}

func (s *AuthorizationServiceImpl) parseTokenToClaims(token string, key []byte) (*model.Claims, error) {
	tkn, err := jwt.ParseWithClaims(token, model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return &model.Claims{}, err
	}
	if tkn == nil || !tkn.Valid {
		return &model.Claims{}, errors.New("Token is incorrect")
	}
	return tkn.Claims.(*model.Claims), nil
}

// GenerateTokens GenerateTokens
func (s *AuthorizationServiceImpl) GenerateTokens(userID uuid.UUID) (accessToken, refreshToken string, accessExp, refreshExp time.Time, err error) {
	accessToken, accessUID, accessExp, err := generateToken(userID, time.Now().Add(s.cfg.AccessTokenTTL), []byte(s.cfg.AccessSecret))
	refreshToken, refreshUID, refreshExp, err := generateToken(userID, time.Now().Add(s.cfg.RefreshTokenTTL), []byte(s.cfg.RefreshSecret))

	if err != nil {
		return
	}

	cashedTokens := model.CachedTokens{
		AccessUID:  accessUID,
		RefreshUID: refreshUID,
	}
	_, err = s.tokenStore.Save(userID.String(), cashedTokens)

	return
}

func generateToken(userID uuid.UUID, expirationTime time.Time, secret []byte) (string, string, time.Time, error) {
	tokenUID := uuid.New()
	claims := &model.Claims{
		UserID: userID,
		UID:    tokenUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", tokenUID.String(), time.Now(), err
	}

	return tokenString, tokenUID.String(), expirationTime, nil
}
