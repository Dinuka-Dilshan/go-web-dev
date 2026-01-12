package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		jwtToken, err := app.authenticator.ValidateToken(parts[1])
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		if jwtToken == nil {
			app.internalServerError(w, r, err)
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			app.unauthorizedError(w, r, err)
			return
		}

		userID, err := strconv.Atoi(fmt.Sprintf("%.f", claims["sub"]))
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		user, err := app.store.Users.GetUserById(r.Context(), userID)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		contextWithUser := context.WithValue(r.Context(), userCtx, user)

		next.ServeHTTP(w, r.WithContext(contextWithUser))

	})
}

func (app *application) checkPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		if user.ID == post.ID {
			next.ServeHTTP(w, r)
			return
		}

		isAllowed, err := app.checkRolePrecedence(r.Context(), role, user)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !isAllowed {
			app.forbiddenError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (app *application) checkRolePrecedence(context context.Context, role string, user *store.User) (bool, error) {
	targetRole, err := app.store.Role.GetByName(context, role)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= targetRole.Level, nil

}
