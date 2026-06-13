package util

import (
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_const "github.com/himbo22/xoxz/account-service/internal/const"
)

type TokenClaims struct {
	UserID   uuid.UUID `json:"uid"`
	RoleID   uuid.UUID `json:"rid"`
	DeviceID string    `json:"did"`
	// Embed standard required claims such as exp, iat, and iss.
	jwt.RegisteredClaims
}

func GenerateAccessToken(
	privateKey *rsa.PrivateKey,
	userID uuid.UUID,
	//roleID uuid.UUID,
	deviceID string,
	ttl time.Duration,
) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		//RoleID:   roleID,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "account-service",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// sign by secret key
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GenerateRefreshToken() string {
	return uuid.NewString()
}

func ParseTokenWithClaims(tokenString string, publicKey *rsa.PublicKey) (*TokenClaims, error) {
	claims := &TokenClaims{}

	// 2. Use ParseWithClaims instead of Parse for typed claims.
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// Prevent algorithm confusion attacks.
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, NewErrorByCode(_const.CodeExpiredAccessToken)
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func HashToken(token string) string {
	hashed := sha256.Sum256([]byte(token))
	return fmt.Sprintf("token:%x", hashed) // fixed 64-char hex
}
