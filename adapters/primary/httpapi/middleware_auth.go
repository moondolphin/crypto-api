package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthContext struct {
	UserID int64
	Email  string
}

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	secret := []byte(jwtSecret)

	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing_token"})
			return
		}

		tokenStr := strings.TrimSpace(h[len("Bearer "):])
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing_token"})
			return
		}

		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			// solo HS256
			if t.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return secret, nil
		})
		if err != nil || tok == nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}

		// sub puede venir como float64 (json) o string
		var uid int64
		switch v := claims["sub"].(type) {
		case float64:
			uid = int64(v)
		case string:
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
				return
			}
			uid = n
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}

		email, _ := claims["email"].(string)
		if uid <= 0 || email == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}

		c.Set("auth", AuthContext{UserID: uid, Email: email})
		c.Next()
	}
}
