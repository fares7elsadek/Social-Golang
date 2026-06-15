package authmiddleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/fares7elsadek/Social-Golang/internal/domain"
	"github.com/fares7elsadek/Social-Golang/internal/services/token"
)

type contextKey string

const claimsKey contextKey = "claims"


func ClaimsFromCtx(ctx context.Context) *domain.Claims {
	v, _ := ctx.Value(claimsKey).(*domain.Claims)
	return v
}


func Authenticate(svc *token.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := r.Header.Get("Authorization")
			if raw == "" {
				writeUnauthorized(w, "missing Authorization header")
				return
			}
 
			parts := strings.SplitN(raw, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeUnauthorized(w, "authorization header must be 'Bearer <token>'")
				return
			}
 
			claims, err := svc.VerifyAccess(parts[1])
			if err != nil {
				switch err {
				case token.ErrTokenExpired:
					writeUnauthorized(w, "token expired")
				default:
					writeUnauthorized(w, "invalid token")
				}
				return
			}
 
			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRoles(roles ...domain.Role) func(http.Handler) http.Handler {
	allowed := make(map[domain.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
 
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromCtx(r.Context())
			if claims == nil {
				// Authenticate middleware must precede RequireRoles.
				writeUnauthorized(w, "not authenticated")
				return
			}
 
			for _, role := range claims.Roles {
				if _, ok := allowed[role]; ok {
					next.ServeHTTP(w, r)
					return
				}
			}
 
			writeForbidden(w)
		})
	}
}


func RequireSelf(userIDFromRequest func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromCtx(r.Context())
			if claims == nil {
				writeUnauthorized(w, "not authenticated")
				return
			}
 
			targetUserID := userIDFromRequest(r)
			if targetUserID == "" || targetUserID != claims.UserID {
				// Allow admins to bypass the self-check.
				for _, role := range claims.Roles {
					if role == domain.RoleAdmin {
						next.ServeHTTP(w, r)
						return
					}
				}
				writeForbidden(w)
				return
			}
 
			next.ServeHTTP(w, r)
		})
	}
}



func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", `Bearer realm="api"`)
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"` + msg + `"}`))
}
 
func writeForbidden(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error":"insufficient permissions"}`))
}