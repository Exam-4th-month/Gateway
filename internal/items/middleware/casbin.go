package middleware

import (
	"gateway-service/internal/items/config"
	"log"

	casbin "github.com/casbin/casbin/v2"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthzMiddleware(path string, enforcer *casbin.Enforcer, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := getRole(c, config)
		ok, err := enforcer.Enforce(role, path, c.Request.Method)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Authorization error"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

func getRole(c *gin.Context, config *config.Config) string {
	tokenString := c.GetHeader("Authorization")
	log.Println(tokenString)
	if tokenString == "" {
		return ""
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(config.JWT.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}

	role, ok := claims["role"].(string)
	if !ok {
		return ""
	}

	return role
}

func GetUser_id(c *gin.Context, config *config.Config) string {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return ""
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(config.JWT.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}

	user_id, ok := claims["user_id"].(string)
	if !ok {
		return ""
	}

	return user_id
}