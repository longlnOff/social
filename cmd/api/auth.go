package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/longlnOff/social/internal/mailer"
	"github.com/longlnOff/social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=250"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token" validate:"required"`
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=250"`
	Password string `json:"password" validate:"required,min=3,max=50"`
}

// createTokenHandler godoc
//
//	@Summary		Create a new token
//	@Description	Creates a new token for the user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken			"Created token"
//	@Failure		500		{object}	string					"Internal Server Error"
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// 1. parse payload
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// 2. validate payload
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// 3. check if user exists or not
	user, err := app.store.User.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.unauthorizedJWTStatelessErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// 4. check if password is correct or not
	if err := user.Password.Check(payload.Password); err != nil {
		app.unauthorizedJWTStatelessErrorResponse(w, r, err)
		return
	}

	// 5. generate the token --> add claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.configuration.Auth.Token.AUTH_TOKEN_EXP).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.configuration.Auth.Token.AUTH_TOKEN_ISS,
		"aud": app.configuration.Auth.Token.AUTH_TOKEN_ISS,
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// 6. return to client
	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// registerUserHandler godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user with the provided username, email, and password
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterUserPayload	true	"User registration data"
//	@Success		201		{object}	UserWithToken		"Registered user"
//	@Failure		400		{object}	string				"Invalid request payload"
//	@Failure		500		{object}	string				"Internal Server Error"
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// hash password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// hash token
	plainToken := uuid.New().String() // plain text token will be attached to email
	hashFunction := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hashFunction[:])

	// store the user
	if err := app.store.User.CreateAndInvite(r.Context(), hashedToken, user, app.configuration.Mail.EXP); err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, store.ErrDuplicateUsername):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	isProduction := app.configuration.Server.ENVIRONMENT == "production"
	activationURL := app.configuration.Server.FRONTEND_URL + "/confirm/" + plainToken
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}
	// send email
	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProduction)
	if err != nil {
		app.logger.Error("Error sending email:", zap.String("error", err.Error()))

		// rollback user creation if email fails (SAGA pattern) - delete user and user's invitation
		if err := app.store.User.Delete(r.Context(), user.ID); err != nil {
			app.logger.Error("Error rolling back user creation:", zap.String("error", err.Error()))
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	app.logger.Info("Email sent to:", zap.String("email", user.Email), zap.Int("status", status))

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
