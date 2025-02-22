// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package entities

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type TeamRole string

const (
	TeamRoleMember TeamRole = "member"
	TeamRoleAdmin  TeamRole = "admin"
)

func (e *TeamRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TeamRole(s)
	case string:
		*e = TeamRole(s)
	default:
		return fmt.Errorf("unsupported scan type for TeamRole: %T", src)
	}
	return nil
}

type NullTeamRole struct {
	TeamRole TeamRole `json:"team_role"`
	Valid    bool     `json:"valid"` // Valid is true if TeamRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTeamRole) Scan(value interface{}) error {
	if value == nil {
		ns.TeamRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TeamRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTeamRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TeamRole), nil
}

type ChatLink struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Hook      int32     `json:"hook"`
	Kind      string    `json:"kind"`
	LinkTo    uuid.UUID `json:"link_to"`
	Data      []byte    `json:"data"`
}

type GithubInstallation struct {
	ID                  uuid.UUID `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	OrgID               uuid.UUID `json:"org_id"`
	InstallationID      int64     `json:"installation_id"`
	InstallationLogin   string    `json:"installation_login"`
	InstallationLoginID int64     `json:"installation_login_id"`
	InstallationType    string    `json:"installation_type"`
	SenderID            int64     `json:"sender_id"`
	SenderLogin         string    `json:"sender_login"`
	IsActive            bool      `json:"is_active"`
}

type GithubOrg struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	OrgID          uuid.UUID `json:"org_id"`
	InstallationID uuid.UUID `json:"installation_id"`
	GithubOrgID    int64     `json:"github_org_id"`
	Name           string    `json:"name"`
}

type GithubRepo struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	InstallationID uuid.UUID `json:"installation_id"`
	GithubID       int64     `json:"github_id"`
	Name           string    `json:"name"`
	FullName       string    `json:"full_name"`
	Url            string    `json:"url"`
	IsActive       bool      `json:"is_active"`
}

type GithubUser struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uuid.UUID `json:"user_id"`
	GithubID    int64     `json:"github_id"`
	GithubOrgID uuid.UUID `json:"github_org_id"`
	Login       string    `json:"login"`
}

type OauthAccount struct {
	ID                uuid.UUID `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	UserID            uuid.UUID `json:"user_id"`
	Provider          string    `json:"provider"`
	ProviderAccountID string    `json:"provider_account_id"`
	ExpiresAt         time.Time `json:"expires_at"`
	Type              string    `json:"type"`
}

type Org struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	Slug      string    `json:"slug"`
	Hooks     []byte    `json:"hooks"`
}

type Repo struct {
	ID            uuid.UUID       `json:"id"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	OrgID         uuid.UUID       `json:"org_id"`
	Name          string          `json:"name"`
	Hook          int32           `json:"hook"`
	HookID        uuid.UUID       `json:"hook_id"`
	DefaultBranch string          `json:"default_branch"`
	IsMonorepo    bool            `json:"is_monorepo"`
	Threshold     int32           `json:"threshold"`
	StaleDuration pgtype.Interval `json:"stale_duration"`
	Url           string          `json:"url"`
	IsActive      bool            `json:"is_active"`
}

type Team struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	OrgID     uuid.UUID `json:"org_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
}

type TeamUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TeamID    uuid.UUID `json:"team_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      TeamRole  `json:"role"`
	IsActive  bool      `json:"is_active"`
	IsAdmin   bool      `json:"is_admin"`
}

type User struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	OrgID      uuid.UUID `json:"org_id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Password   string    `json:"password"`
	Picture    string    `json:"picture"`
	IsActive   bool      `json:"is_active"`
	IsVerified bool      `json:"is_verified"`
}

type UserRole struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	UserID    uuid.UUID `json:"user_id"`
	OrgID     uuid.UUID `json:"org_id"`
}
