// Package github provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/breuHQ/oapi-codegen, a modified copy of github.com/deepmap/oapi-codegen/v2.
//
// It was modified to add support for the following features:
//  - Support for custom templates by filename.
//  - Supporting x-breu-entity in the schema to generate a struct for the entity.
//
// DO NOT EDIT!!

package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	"github.com/scylladb/gocqlx/v2/table"
	"go.breu.io/quantm/internal/shared"
)

const (
	APIKeyAuthScopes = "APIKeyAuth.Scopes"
	BearerAuthScopes = "BearerAuth.Scopes"
)

var (
	ErrInvalidSetupAction    = errors.New("invalid SetupAction value")
	ErrInvalidWorkflowStatus = errors.New("invalid WorkflowStatus value")
)

type (
	SetupActionMapType map[string]SetupAction // SetupActionMapType is a quick lookup map for SetupAction.
)

// Defines values for SetupAction.
const (
	SetupActionDelete                 SetupAction = "delete"
	SetupActionInstall                SetupAction = "install"
	SetupActionNewPermissionsAccepted SetupAction = "new_permissions_accepted"
	SetupActionRequest                SetupAction = "request"
	SetupActionSuspend                SetupAction = "suspend"
	SetupActionUnsuspend              SetupAction = "unsuspend"
	SetupActionUpdate                 SetupAction = "update"
)

// SetupActionMap returns all known values for SetupAction.
var (
	SetupActionMap = SetupActionMapType{
		SetupActionDelete.String():                 SetupActionDelete,
		SetupActionInstall.String():                SetupActionInstall,
		SetupActionNewPermissionsAccepted.String(): SetupActionNewPermissionsAccepted,
		SetupActionRequest.String():                SetupActionRequest,
		SetupActionSuspend.String():                SetupActionSuspend,
		SetupActionUnsuspend.String():              SetupActionUnsuspend,
		SetupActionUpdate.String():                 SetupActionUpdate,
	}
)

/*
 * Helper methods for SetupAction for easy marshalling and unmarshalling.
 */
func (v SetupAction) String() string               { return string(v) }
func (v SetupAction) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *SetupAction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := SetupActionMap[s]
	if !ok {
		return ErrInvalidSetupAction
	}

	*v = val

	return nil
}

type (
	WorkflowStatusMapType map[string]WorkflowStatus // WorkflowStatusMapType is a quick lookup map for WorkflowStatus.
)

// Defines values for WorkflowStatus.
const (
	WorkflowStatusFailure  WorkflowStatus = "failure"
	WorkflowStatusQueued   WorkflowStatus = "queued"
	WorkflowStatusSignaled WorkflowStatus = "signaled"
	WorkflowStatusSkipped  WorkflowStatus = "skipped"
	WorkflowStatusSuccess  WorkflowStatus = "success"
)

// WorkflowStatusMap returns all known values for WorkflowStatus.
var (
	WorkflowStatusMap = WorkflowStatusMapType{
		WorkflowStatusFailure.String():  WorkflowStatusFailure,
		WorkflowStatusQueued.String():   WorkflowStatusQueued,
		WorkflowStatusSignaled.String(): WorkflowStatusSignaled,
		WorkflowStatusSkipped.String():  WorkflowStatusSkipped,
		WorkflowStatusSuccess.String():  WorkflowStatusSuccess,
	}
)

/*
 * Helper methods for WorkflowStatus for easy marshalling and unmarshalling.
 */
func (v WorkflowStatus) String() string               { return string(v) }
func (v WorkflowStatus) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }
func (v *WorkflowStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, ok := WorkflowStatusMap[s]
	if !ok {
		return ErrInvalidWorkflowStatus
	}

	*v = val

	return nil
}

// CompleteInstallationRequest complete the installation given the installation_id & setup_action.
type CompleteInstallationRequest struct {
	InstallationID shared.Int64 `json:"installation_id"`
	SetupAction    SetupAction  `json:"setup_action"`
}

// CreateGithubUserOrgsRequest defines model for CreateGithubUserOrgsRequest.
type CreateGithubUserOrgsRequest struct {
	GithubOrgIDs []shared.Int64 `json:"github_org_ids"`
	GithubUserID shared.Int64   `json:"github_user_id"`
	UserID       gocql.UUID     `json:"user_id"`
}

// CreateTeamUserRequest defines model for CreateTeamUserRequest.
type CreateTeamUserRequest struct {
	GithubOrgID  shared.Int64 `json:"github_org_id"`
	GithubUserID shared.Int64 `json:"github_user_id"`
	TeamID       gocql.UUID   `json:"team_id"`
	UserID       gocql.UUID   `json:"user_id"`
}

// Installation defines model for GithubInstallation.
type Installation struct {
	CreatedAt           time.Time    `json:"created_at"`
	ID                  gocql.UUID   `json:"id"`
	InstallationID      shared.Int64 `json:"installation_id" validate:"required,db_unique"`
	InstallationLogin   string       `json:"installation_login"`
	InstallationLoginID shared.Int64 `json:"installation_login_id"`
	InstallationType    string       `json:"installation_type"`
	SenderID            shared.Int64 `json:"sender_id"`
	SenderLogin         string       `json:"sender_login"`
	Status              string       `json:"status"`
	TeamID              gocql.UUID   `json:"team_id"`
	UpdatedAt           time.Time    `json:"updated_at"`
}

var (
	githubinstallationMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "github_installations",
			Columns: []string{"created_at", "id", "installation_id", "installation_login", "installation_login_id", "installation_type", "sender_id", "sender_login", "status", "team_id", "updated_at"},
			PartKey: []string{"id", "team_id"},
		},
	}

	githubinstallationTable = itable.New(*githubinstallationMeta.M)
)

func (githubinstallation *Installation) GetTable() itable.ITable {
	return githubinstallationTable
}

// OrgUser defines model for GithubOrgUser.
type OrgUser struct {
	CreatedAt     time.Time    `json:"created_at"`
	GithubOrgID   shared.Int64 `json:"github_org_id"`
	GithubOrgName string       `json:"github_org_name"`
	GithubUserID  shared.Int64 `json:"github_user_id"`
	ID            gocql.UUID   `json:"id"`
	UpdatedAt     time.Time    `json:"updated_at"`

	// UserId auth's user ID.
	UserID gocql.UUID `json:"user_id"`
}

var (
	githuborguserMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "github_org_users",
			Columns: []string{"created_at", "github_org_id", "github_org_name", "github_user_id", "id", "updated_at", "user_id"},
			PartKey: []string{"id", "user_id"},
		},
	}

	githuborguserTable = itable.New(*githuborguserMeta.M)
)

func (githuborguser *OrgUser) GetTable() itable.ITable {
	return githuborguserTable
}

// Repo defines model for GithubRepo.
type Repo struct {
	CreatedAt       time.Time    `json:"created_at"`
	DefaultBranch   string       `json:"default_branch"`
	FullName        string       `json:"full_name"`
	GithubID        shared.Int64 `json:"github_id"`
	HasEarlyWarning bool         `json:"has_early_warning"`
	ID              gocql.UUID   `json:"id"`
	InstallationID  shared.Int64 `json:"installation_id"`
	IsActive        bool         `json:"is_active"`
	Name            string       `json:"name"`
	TeamID          gocql.UUID   `json:"team_id"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

var (
	githubrepoMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "github_repos",
			Columns: []string{"created_at", "default_branch", "full_name", "github_id", "has_early_warning", "id", "installation_id", "is_active", "name", "team_id", "updated_at"},
			PartKey: []string{"id", "team_id"},
		},
	}

	githubrepoTable = itable.New(*githubrepoMeta.M)
)

func (githubrepo *Repo) GetTable() itable.ITable {
	return githubrepoTable
}

// SetupAction defines model for SetupAction.
type SetupAction string

// WorkflowResponse workflow status & run id.
type WorkflowResponse struct {
	RunID string `json:"run_id"`

	// Status the workflow status
	Status WorkflowStatus `json:"status"`
}

// WorkflowStatus the workflow status.
type WorkflowStatus string

// GithubGetInstallationsParams defines parameters for GithubGetInstallations.
type GithubGetInstallationsParams struct {
	// InstallationLogin github team name
	InstallationLogin *string `form:"installation_login,omitempty" json:"installation_login,omitempty"`

	// InstallationId installation ID of the github app.
	InstallationId *shared.Int64 `form:"installation_id,omitempty" json:"installation_id,omitempty"`
}

// GithubListUserOrgsParams defines parameters for GithubListUserOrgs.
type GithubListUserOrgsParams struct {
	// UserId User ID
	UserId string `form:"user_id" json:"user_id"`
}

// GithubCompleteInstallationJSONRequestBody defines body for GithubCompleteInstallation for application/json ContentType.
type GithubCompleteInstallationJSONRequestBody = CompleteInstallationRequest

// GithubCreateUserOrgsJSONRequestBody defines body for GithubCreateUserOrgs for application/json ContentType.
type GithubCreateUserOrgsJSONRequestBody = CreateGithubUserOrgsRequest

// CreateTeamUserJSONRequestBody defines body for CreateTeamUser for application/json ContentType.
type CreateTeamUserJSONRequestBody = CreateTeamUserRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Complete GitHub App installation
	// (POST /providers/github/complete-installation)
	GithubCompleteInstallation(ctx echo.Context) error

	// Get GitHub installations
	// (GET /providers/github/installations)
	GithubGetInstallations(ctx echo.Context, params GithubGetInstallationsParams) error

	// Get GitHub repositories
	// (GET /providers/github/repos)
	GithubGetRepos(ctx echo.Context) error

	// list assoicated organizations
	// (GET /providers/github/user-orgs)
	GithubListUserOrgs(ctx echo.Context, params GithubListUserOrgsParams) error

	// associate github organizations for the newly registered user.
	// (POST /providers/github/user-orgs)
	GithubCreateUserOrgs(ctx echo.Context) error

	// Webhook reciever for github
	// (POST /providers/github/webhook)
	GithubWebhook(ctx echo.Context) error

	// associate team and team users for the newly registered user.
	// (POST /providers/teams/team-user)
	CreateTeamUser(ctx echo.Context) error

	// SecurityHandler returns the underlying Security Wrapper
	SecureHandler(ctx echo.Context, handler echo.HandlerFunc) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GithubCompleteInstallation converts echo context to params.

func (w *ServerInterfaceWrapper) GithubCompleteInstallation(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	ctx.Set(APIKeyAuthScopes, []string{})

	handler := func(ctx echo.Context) error {
		return w.Handler.GithubCompleteInstallation(ctx)
	}
	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SecureHandler(ctx, handler)

	return err
}

// GithubGetInstallations converts echo context to params.

func (w *ServerInterfaceWrapper) GithubGetInstallations(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GithubGetInstallationsParams
	// ------------- Optional query parameter "installation_login" -------------

	err = runtime.BindQueryParameter("form", true, false, "installation_login", ctx.QueryParams(), &params.InstallationLogin)
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, fmt.Errorf("Invalid format for parameter installation_login: %s", err))
	}

	// ------------- Optional query parameter "installation_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "installation_id", ctx.QueryParams(), &params.InstallationId)
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, fmt.Errorf("Invalid format for parameter installation_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GithubGetInstallations(ctx, params)

	return err
}

// GithubGetRepos converts echo context to params.

func (w *ServerInterfaceWrapper) GithubGetRepos(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	ctx.Set(APIKeyAuthScopes, []string{})

	handler := func(ctx echo.Context) error {
		return w.Handler.GithubGetRepos(ctx)
	}
	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SecureHandler(ctx, handler)

	return err
}

// GithubListUserOrgs converts echo context to params.

func (w *ServerInterfaceWrapper) GithubListUserOrgs(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GithubListUserOrgsParams
	// ------------- Required query parameter "user_id" -------------

	err = runtime.BindQueryParameter("form", true, true, "user_id", ctx.QueryParams(), &params.UserId)
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, fmt.Errorf("Invalid format for parameter user_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GithubListUserOrgs(ctx, params)

	return err
}

// GithubCreateUserOrgs converts echo context to params.

func (w *ServerInterfaceWrapper) GithubCreateUserOrgs(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GithubCreateUserOrgs(ctx)

	return err
}

// GithubWebhook converts echo context to params.

func (w *ServerInterfaceWrapper) GithubWebhook(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GithubWebhook(ctx)

	return err
}

// CreateTeamUser converts echo context to params.

func (w *ServerInterfaceWrapper) CreateTeamUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateTeamUser(ctx)

	return err
}

// EchoRouter is an interface that wraps the methods of echo.Echo & echo.Group to provide a common interface
// for registering routes.
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/providers/github/complete-installation", wrapper.GithubCompleteInstallation)
	router.GET(baseURL+"/providers/github/installations", wrapper.GithubGetInstallations)
	router.GET(baseURL+"/providers/github/repos", wrapper.GithubGetRepos)
	router.GET(baseURL+"/providers/github/user-orgs", wrapper.GithubListUserOrgs)
	router.POST(baseURL+"/providers/github/user-orgs", wrapper.GithubCreateUserOrgs)
	router.POST(baseURL+"/providers/github/webhook", wrapper.GithubWebhook)
	router.POST(baseURL+"/providers/teams/team-user", wrapper.CreateTeamUser)

}
