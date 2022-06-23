package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtWrapper struct {
	SecretKey       string
	ExpirationHours int64
}

type JwtClaim struct {
	Name string
	jwt.StandardClaims
}

func NewJwtWrapper(secretKey string, expirationHours int64) *JwtWrapper {
	return &JwtWrapper{
		SecretKey:       secretKey,
		ExpirationHours: expirationHours,
	}
}

func (j *JwtWrapper) GenerateToken(name string) (signedToken string, err error) {
	claims := &JwtClaim{
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}

	return signedToken, nil
}

func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")

		return
	}

	return
}

func (j *JwtWrapper) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Missing Authorization Header"))
			if err != nil {
				log.Fatal(err)
			}

			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		_, err := j.ValidateToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Error verifying JWT token: " + err.Error()))
			if err != nil {
				log.Fatal(err)
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}
