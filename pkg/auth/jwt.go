package auth

import (
	"errors"
	"globalZT/pkg/config"

	"github.com/dgrijalva/jwt-go"
)

var jwtTools JwtTools

type JwtTools struct {
	SigningKey  []byte
	ExpiresTime int64
	BufferTime  int64
	Issue       string
}

func NewJwtTools(c config.AUTH) {
	jwtTools = JwtTools{
		SigningKey:  []byte(c.SigningKey),
		ExpiresTime: c.ExpiresTime,
		BufferTime:  c.BufferTime,
		Issue:       "globalZT-center",
	}
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

func (j *JwtTools) CreateToken(claims JwtInfo) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JwtTools) ParseToken(tokenString string) (*JwtInfo, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtInfo{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*JwtInfo); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}
}

type JwtInfo struct {
	UID        uint
	GID        uint
	RID        uint
	User       string
	BufferTime int64
	token      string
	jwt.StandardClaims
}
