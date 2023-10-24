package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

//Claims that passing via jwt token
type Claims struct {
	Email   string   		`json:"email"`
	Role    string          `json:"role"`
	jwt.StandardClaims
}

//GenerateToken to generate jwt token
func GenerateToken(userEmail string, userRole string) (string,error){
	claims:=Claims{
		userEmail,
		userRole,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour*24).Unix(),
			IssuedAt: time.Now().Unix(),
		},
	}

	token :=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	secretKey := []byte("101101")

	tokenString,err :=token.SignedString(secretKey)
	if err != nil{
		return "",err
	}

	return tokenString,nil
}

