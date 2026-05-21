package auth

import (
	"os"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const Issuer = "token"

func Jwt(userID string) (string, error) {
	jwtSecret := os.Getenv("SECRET_JWT")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    Issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour)),
		Subject:   userID,
	})

	return token.SignedString([]byte(jwtSecret))
}

func ValidateJwt(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := token.Claims.GetSubject()
	if userID == "" || err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil || issuer != Issuer {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func HashPassword(password string) (string, error) {
	hashpassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hashpassword, nil
}

func CheckPassword(password, hashPassword string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hashPassword)
	if err != nil || !match {
		return false, err
	}

	return true, nil
}

func ParseUUID(s string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}
