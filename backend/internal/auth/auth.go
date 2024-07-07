package auth

import (
	"backend/internal/config"
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var cfg, _ = config.LoadConfig() // Consider handling this error properly

type CustomClaims struct {
	Email_ID string `json:"emailid"`
	jwt.StandardClaims
}

func GenerateTokens(userID string) (accessToken string, refreshToken string, err error) {
	// Generate access token
	t := time.Now().Local().Add(time.Minute * 5).Unix()
	log.Println("time assigned ", t)
	claims := CustomClaims{
		Email_ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: t,
			Issuer:    userID,
		},
	}
	accessTokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessTokenClaim.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	t = time.Now().Local().Add(time.Minute * 15).Unix()
	refreshTokenClaims := CustomClaims{
		Email_ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: t,
			Issuer:    userID,
		},
	}
	refreshTokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshTokenClaim.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func RefreshAccessToken(refreshTokenString string) (newAccessToken string, err error) {
	// Validate the refresh token
	token, err := jwt.ParseWithClaims(refreshTokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", jwt.ErrInvalidKey
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return "", errors.New("invalid claims token")
	}

	// Generate a new access token
	newAccessToken, _, err = GenerateTokens(claims.Email_ID)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
