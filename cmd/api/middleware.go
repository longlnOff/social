package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/longlnOff/social/internal/store"
)

// AuthorizationMiddleware for post
func (app *application) checkPostownership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get user from context
		user := getUserFromCtx(r)

		// 2. Get post from context
		post := getPostFromCtx(r)

		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		// 3. role precedence check
		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenErrorResponse(w, r, fmt.Errorf("permission denied"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, requiredRole string) (bool, error) {
	role, err := app.store.Role.GetByName(ctx, requiredRole)
	if err != nil {
		return false, err
	}

	return role.Level <= user.Role.Level, nil
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. read the auth header
		// Add this right before checking the auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedJWTStatelessErrorResponse(w, r, fmt.Errorf("no auth header"))
			return
		}
		// 2. parse it --> get the token
		parts := strings.Split(authHeader, " ") // Authorization: Bearer <token>
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedJWTStatelessErrorResponse(w, r, fmt.Errorf("invalid auth header"))
			return
		}
		token := parts[1]

		// 3. validate the token
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedJWTStatelessErrorResponse(w, r, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedJWTStatelessErrorResponse(w, r, err)
			return
		}
		ctx := r.Context()
		// 4. add the token to the context
		user, err := app.store.User.GetByUserID(ctx, userID)
		if err != nil {
			app.unauthorizedJWTStatelessErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, Userctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicAuthErrorResponse(w, r, fmt.Errorf("no auth header"))
				return
			}
			// parse it --> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicAuthErrorResponse(w, r, fmt.Errorf("invalid auth header"))
				return
			}
			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicAuthErrorResponse(w, r, fmt.Errorf("invalid auth credentials"))
				return
			}
			// check the credentials
			credentials := strings.SplitN(string(decoded), ":", 2)
			if len(credentials) != 2 {
				app.unauthorizedBasicAuthErrorResponse(w, r, fmt.Errorf("invalid auth credentials"))
				return
			}
			account, password := credentials[0], credentials[1]

			username := app.configuration.Auth.Basic.AUTH_BASIC_USER
			pass := app.configuration.Auth.Basic.AUTH_BASIC_PASSWORD
			if account != username || password != pass {
				app.unauthorizedBasicAuthErrorResponse(w, r, fmt.Errorf("invalid auth credentials"))
				return
			}

			// call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
