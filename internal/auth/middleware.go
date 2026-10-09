package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware intercepta las peticiones y valida el token Bearer en la cabecera Authorization
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Cabecera Authorization requerida"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error":"Formato de token inválido. Usar 'Bearer <token>'"}`, http.StatusUnauthorized)
			return
		}

		userID, _, err := ValidateJWT(parts[1])
		if err != nil {
			http.Error(w, `{"error":"Token inválido o expirado"}`, http.StatusUnauthorized)
			return
		}

		// Inyectar el userID en el contexto de la petición HTTP
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
