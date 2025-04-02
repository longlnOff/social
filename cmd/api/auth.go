package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"net/http"

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

	// send email

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
