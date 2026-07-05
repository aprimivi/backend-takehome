package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type TokenClaims struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"expires_at"`
}

func GenerateToken(userID int64, email string) (string, error) {
	claims := TokenClaims{
		UserID:    userID,
		Email:     email,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := signToken(encodedPayload)
	return fmt.Sprintf("%s.%s", encodedPayload, signature), nil
}

func ValidateToken(token string) (TokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return TokenClaims{}, errors.New("invalid token")
	}

	expectedSignature := signToken(parts[0])
	if !hmac.Equal([]byte(parts[1]), []byte(expectedSignature)) {
		return TokenClaims{}, errors.New("invalid token signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return TokenClaims{}, err
	}

	var claims TokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return TokenClaims{}, err
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return TokenClaims{}, errors.New("token expired")
	}

	return claims, nil
}

func signToken(payload string) string {
	secret := []byte(getSecret())
	m := hmac.New(sha256.New, secret)
	_, _ = m.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

func getSecret() string {
	if value := os.Getenv("APP_TOKEN_SECRET"); value != "" {
		return value
	}

	return "dev-token-secret"
}
