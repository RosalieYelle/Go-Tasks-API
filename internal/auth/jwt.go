package auth

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "net/http"
    "time"
)

// Normally  not hard coded.
var jwtKey = []byte("secret-key")

func GenerateToken(userId string) (string, error) {
    claims := jwt.MapClaims{
        "userId": userId,
		//need to log in again after one day
        "exp":    time.Now().Add(time.Hour * 24).Unix(),
    }
	//token with HMAC-SHA256
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
            return
        }
		// extract token, remove first 6 char (Bearer) Provides the secret key via the callback 
        token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }
		// access custom claims
        claims := token.Claims.(jwt.MapClaims)
		//saves user id
        c.Set("userId", claims["userId"].(string))
		//proceed to next handler
        c.Next()
    }
