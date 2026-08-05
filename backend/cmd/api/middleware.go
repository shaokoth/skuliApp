package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/skuliapp/backend/internal/users"
)

// recoverPanic converts a panic in a handler into a clean 500 response.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// logRequest logs one line per request.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("request",
			"ip", r.RemoteAddr,
			"proto", r.Proto,
			"method", r.Method,
			"uri", r.URL.RequestURI(),
		)
		next.ServeHTTP(w, r)
	})
}

// requireAuth verifies the Bearer token, loads the identity into the request
// context and rejects the request when the token is missing or invalid.
func (app *application) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.authenticationRequiredResponse(w, r)
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		claims, err := app.jwt.Parse(parts[1])
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		r = app.contextSetUser(r, authUser{
			UserID:   claims.UserID,
			SchoolID: claims.SchoolID,
			Role:     claims.Role,
		})
		next(w, r)
	}
}

// requireRole builds on requireAuth and additionally checks the caller's role.
func (app *application) requireRole(next http.HandlerFunc, roles ...users.Role) http.HandlerFunc {
	return app.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		u, ok := app.contextGetUser(r)
		if !ok {
			app.authenticationRequiredResponse(w, r)
			return
		}

		permitted := false
		for _, role := range roles {
			if string(role) == u.Role {
				permitted = true
				break
			}
		}
		if !permitted {
			app.notPermittedResponse(w, r)
			return
		}
		next(w, r)
	})
}
