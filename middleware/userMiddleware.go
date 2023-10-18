package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)


func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString==""{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Token is missing"})
			c.Abort()
			return
		}

		tokenString= strings.Replace(tokenString, "Bearer ", "", 1)

		token,err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{},error) {
			return []byte("101101"),nil
		})

		if err != nil || !token.Valid{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid token"})
			c.Abort()
			return
		}

		claims,ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid token claims"})
			c.Abort()
			return
		}

		email,ok := claims["email"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Email missing in claims"})
			c.Abort()
			return
		}
		c.Set("email",email)

		c.Next()
	}
}