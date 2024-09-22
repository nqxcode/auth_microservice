package config

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

const (
	authRefreshTokenSecretKeyEnvName  = "AUTH_REFRESH_TOKEN_SECRET"
	authAccessTokenSecretKeyEnvName   = "AUTH_ACCESS_TOKEN"
	authRefreshTokenExpirationEnvName = "AUTH_REFRESH_TOKEN_EXPIRATION"
	authAccessTokenExpirationEnvName  = "AUTH_ACCESS_TOKEN_EXPIRATION"
)

const (
	defaultRefreshTokenExpiration = 60 * time.Minute
	defaultAccessTokenExpiration  = 5 * time.Minute
)

// AuthConfig auth config
type AuthConfig interface {
	RefreshTokenSecretKey() string
	AccessTokenSecretKey() string
	RefreshTokenExpiration() time.Duration
	AccessTokenExpiration() time.Duration
}

type authConfig struct {
	refreshTokenSecretKey  string
	accessTokenSecretKey   string
	refreshTokenExpiration time.Duration
	accessTokenExpiration  time.Duration
}

func NewAuthConfig() (AuthConfig, error) {
	refreshTokenSecretKey := os.Getenv(authRefreshTokenSecretKeyEnvName)
	if len(refreshTokenSecretKey) == 0 {
		return nil, errors.New("refresh token secret key not found")
	}

	accessTokenSecretKey := os.Getenv(authAccessTokenSecretKeyEnvName)
	if len(accessTokenSecretKey) == 0 {
		return nil, errors.New("access token secret key not found")
	}

	refreshTokenExpirationStr := os.Getenv(authRefreshTokenExpirationEnvName)
	if len(refreshTokenExpirationStr) == 0 {
		refreshTokenExpirationStr = strconv.Itoa(int(defaultRefreshTokenExpiration.Minutes()))
	}
	refreshTokenExpiration, err := strconv.Atoi(refreshTokenExpirationStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse refresh token expiration")
	}

	accessTokenExpirationStr := os.Getenv(authAccessTokenExpirationEnvName)
	if len(accessTokenExpirationStr) == 0 {
		accessTokenExpirationStr = strconv.Itoa(int(defaultAccessTokenExpiration.Minutes()))
	}
	accessTokenExpiration, err := strconv.Atoi(accessTokenExpirationStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse access token expiration")
	}

	return &authConfig{
		refreshTokenSecretKey:  refreshTokenSecretKey,
		accessTokenSecretKey:   accessTokenSecretKey,
		refreshTokenExpiration: time.Duration(refreshTokenExpiration) * time.Minute,
		accessTokenExpiration:  time.Duration(accessTokenExpiration) * time.Minute,
	}, nil
}

func (c *authConfig) RefreshTokenSecretKey() string {
	return c.refreshTokenSecretKey
}

func (c *authConfig) AccessTokenSecretKey() string {
	return c.accessTokenSecretKey
}

func (c *authConfig) RefreshTokenExpiration() time.Duration {
	return c.refreshTokenExpiration
}

func (c *authConfig) AccessTokenExpiration() time.Duration {
	return c.accessTokenExpiration
}
