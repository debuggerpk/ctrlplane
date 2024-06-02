// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	// SecurityHandler is the base security handler for the API. It is meant to be embedded in other handlers.
	//
	// Usage:
	//  package {name}
	//
	//  import "go.breu.io/quantm/internal/auth"
	//
	//  type ServerHandler struct {
	//    *auth.SecurityHandler
	//  }
	SecurityHandler struct{ Middleware echo.MiddlewareFunc }
	ServerHandler   struct{ *SecurityHandler } // ServerHandler for auth
)

// NewServerHandler creates a new ServerHandler.
func NewServerHandler(security echo.MiddlewareFunc) *ServerHandler {
	return &ServerHandler{
		SecurityHandler: &SecurityHandler{Middleware: security},
	}
}

// SecureHandler wraps the handler with the security middleware.
func (s *SecurityHandler) SecureHandler(ctx echo.Context, handler echo.HandlerFunc) error {
	err := s.Middleware(handler)(ctx)
	return err
}

// endpoint: /auth/register
func (s *ServerHandler) Register(ctx echo.Context) error {
	request := &RegisterationRequest{}

	// Translating request to json
	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	// Validating request
	if err := ctx.Validate(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	// Validating team
	team := &Team{Name: request.TeamName}
	if err := ctx.Validate(team); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	user := &User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}
	if err := ctx.Validate(user); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if err := db.Save(team); err != nil {
		// TODO: cleanup created team.
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	user.TeamID = team.ID
	if err := db.Save(user); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, &RegisterationResponse{Team: team, User: user})
}

// endpoint: /auth/login
func (s *ServerHandler) Login(ctx echo.Context) error {
	request := &LoginRequest{}

	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	params := db.QueryParams{"email": "'" + request.Email + "'"}
	user := &User{}

	if err := db.Get(user, params); err != nil {
		return shared.NewAPIError(http.StatusNotFound, ErrInvalidCredentials)
	}

	if user.VerifyPassword(request.Password) {
		access, _ := GenerateAccessToken(user.ID.String(), user.TeamID.String())
		refresh, _ := GenerateRefreshToken(user.ID.String(), user.TeamID.String())

		return ctx.JSON(http.StatusOK, &TokenResponse{AccessToken: access, RefreshToken: refresh})
	}

	return shared.NewAPIError(http.StatusUnauthorized, ErrInvalidCredentials)
}

// endpoint: auth/api-keys/user
func (s *ServerHandler) CreateUserAPIKey(ctx echo.Context) error {
	request := &CreateAPIKeyRequest{}

	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	id, _ := gocql.ParseUUID(ctx.Get("user_id").(string))
	guard := &Guard{}
	key := guard.NewForUser(*request.Name, id)

	if err := guard.Save(); err != nil {
		shared.Logger().Error("error saving guard", "error", err)
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, &CreateAPIKeyResponse{Key: &key})
}

// endpoint: /auth/api-keys/team
func (s *ServerHandler) CreateTeamAPIKey(ctx echo.Context) error {
	request := &CreateAPIKeyRequest{}

	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	id, _ := gocql.ParseUUID(ctx.Get("team_id").(string))
	guard := &Guard{}
	key := guard.NewForTeam(id)

	if err := guard.Save(); err != nil {
		shared.Logger().Error("error saving guard", "error", err)
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, &CreateAPIKeyResponse{Key: &key})
}

// endpoint: /auth/api-keys/validate
func (s *ServerHandler) ValidateAPIKey(ctx echo.Context) error {
	valid := "valid" // FIXME: this is not correct.
	return ctx.JSON(http.StatusOK, &APIKeyValidationResponse{Message: &valid})
}

// endpoint: /auth/teams
func (s *ServerHandler) CreateTeam(ctx echo.Context) error {
	request := &CreateTeamRequest{}

	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	id := ctx.Get("user_id").(string)

	user := &User{}
	team := &Team{Name: request.Name}

	if err := db.Get(user, db.QueryParams{"id": id}); err != nil {
		slog.Error("error getting user", "error", err)
		return shared.NewAPIError(http.StatusNotFound, err)
	}

	if err := db.Save(team); err != nil {
		slog.Error("error saving team", "error", err)
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	user.TeamID = team.ID
	if err := db.Save(user); err != nil {
		slog.Error("error creating user", "error", err)
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, team)
}

// endpoint: /auth/teams
func (s *ServerHandler) ListTeams(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// endpoint: /auth/teams/:id
func (s *ServerHandler) GetTeam(ctx echo.Context, id gocql.UUID) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

func (s *ServerHandler) SetActiveTeam(ctx echo.Context, id gocql.UUID) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// endpoint: /auth/teams/:id/users
func (s *ServerHandler) AddUserToTeam(ctx echo.Context, id gocql.UUID) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// endpoint: /auth/users
func (s *ServerHandler) CreateUser(ctx echo.Context) error {
	request := &CreateUserRequest{}

	// Translating request to json
	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	// Validating request
	if err := ctx.Validate(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	names := strings.Split(request.Name, " ")

	user := &User{
		FirstName:  names[0],
		LastName:   "", // Default value
		Email:      request.Email,
		Password:   "", // Default value
		IsActive:   true,
		IsVerified: true,
	}

	if request.TeamID != nil {
		user.TeamID = *request.TeamID
	}

	// Check if names slice has at least 2 elements, if so, assign the second element to LastName
	if len(names) >= 2 {
		user.LastName = strings.Join(names[1:], " ")
	}

	if err := ctx.Validate(user); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if err := db.Save(user); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, user)
}

// ListUsers handles the following use cases
//
//   - List all users associated with the team
//   - Get a user if an email is given in query params
//   - Get a user if social account id is given in query params
//
// endpoint: /auth/users
func (s *ServerHandler) ListUsers(ctx echo.Context, params ListUsersParams) error {
	users := make([]User, 0)

	if params.Email != nil {
		if err := db.Filter(&User{}, users, db.QueryParams{"email": *params.Email}); err != nil {
			return shared.NewAPIError(http.StatusBadRequest, err)
		}
	}

	if params.ProviderAccountId != nil && params.Provider != nil {
		filter := db.QueryParams{"provider": quote(*params.Provider), "provider_account_id": quote(*params.ProviderAccountId)}
		if err := db.Filter(&Account{}, users, filter); err != nil {
			return shared.NewAPIError(http.StatusBadRequest, err)
		}
	}

	return ctx.JSON(http.StatusOK, users)
}

// endpoint: /auth/users/:id
func (s *ServerHandler) GetUser(ctx echo.Context, id string) error {
	user := &User{}
	param := db.QueryParams{"id": ctx.Param("id")}

	if err := db.Get(user, param); err != nil {
		return shared.NewAPIError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, user)
}

// endpoint: /auth/users/:id
func (s *ServerHandler) UpdateUser(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

// endpoint: /auth/accounts/link
// TODO: should be /auth/users/:id/link-account.
func (s *ServerHandler) LinkAccount(ctx echo.Context) error {
	request := &LinkAccountRequest{}
	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	account := &Account{
		UserID:            request.UserID,
		Provider:          request.Provider,
		ProviderAccountID: request.ProviderAccountID,
		ExpiresAt:         request.ExpiresAt,
		Type:              request.Type,
	}

	if err := db.Save(account); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, account)
}
