package middleware

import (
	"net/http"
	"os"
	"strings"
	"vaccination-server/models"

	"github.com/golang-jwt/jwt"
)

var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

func shoulCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

func CheckAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !shoulCheckToken(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})

}
