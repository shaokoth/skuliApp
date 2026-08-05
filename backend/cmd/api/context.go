package main

import (
	"context"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

// authUser is the authenticated identity extracted from a valid JWT.
type authUser struct {
	UserID   int64
	SchoolID int64
	Role     string
}

func (app *application) contextSetUser(r *http.Request, u authUser) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, u)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) (authUser, bool) {
	u, ok := r.Context().Value(userContextKey).(authUser)
	return u, ok
}
