package helper

import (
	"time"

	"github.com/RuhullahReza/SecondHand/util/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Payload struct {
	UserId uuid.UUID `json:"user_id"`
	Name string `json:"username"`
	Role string `json:"role"`
}

type Claims struct {
	Payload Payload
	jwt.RegisteredClaims
}

var jwtKey = []byte(config.JwtKey())

func SignToken(id uuid.UUID, name string,role string) (string,error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Payload: Payload{
			UserId :id,
			Name : name,
			Role : role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token :=  jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return tokenString,err
	}
	return tokenString, nil
}

func ValidateToken(tokenStr string)(Payload, bool){
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid  {
		return claims.Payload, false
	}

	return claims.Payload, true
}