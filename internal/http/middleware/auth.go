package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"golang-htmx-bulma/internal/service"
)

func AuthMiddleware(signingKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			// If it's an HTMX request, we might want to return a special header to redirect
			if c.GetHeader("HX-Request") == "true" {
				c.Header("HX-Redirect", "/login")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// Parse and validate token
		claims := &service.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(signingKey), nil
		})

		if err != nil || !token.Valid {
			// Clear invalid cookie
			c.SetCookie("jwt_token", "", -1, "/", "", false, true)
			
			if c.GetHeader("HX-Request") == "true" {
				c.Header("HX-Redirect", "/login")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// Set claims in context
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.RoleID)
		c.Set("user_first_name", claims.FirstName)
		c.Set("user_last_name", claims.LastName)

		c.Next()
	}
}
