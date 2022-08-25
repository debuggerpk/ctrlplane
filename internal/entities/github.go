package entities

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
)

var (
	githubInstallationColumns = []string{
		"id",
		"team_id",
		"installation_id",
		"installation_login",
		"installation_type",
		"sender_id",
		"sender_login",
		"status",
		"created_at",
		"updated_at",
	}

	githubInstallationMeta = table.Metadata{
		Name:    "github_installations",
		Columns: githubInstallationColumns,
		PartKey: []string{"id"},
	}

	githubInstallationTable = table.New(githubInstallationMeta)
)

type GithubInstallation struct {
	ID                gocql.UUID `json:"id" cql:"id"`
	TeamID            gocql.UUID `json:"team_id" cql:"team_id"`
	InstallationID    int64      `json:"installation_id" cql:"installation_id" validate:"required,db_unique"`
	InstallationLogin string     `json:"installation_login" cql:"installation_login"`
	InstallationType  string     `json:"installation_type" cql:"installation_type"`
	SenderID          int64      `json:"sender_id" cql:"sender_id"`
	SenderLogin       string     `json:"sender_login" cql:"sender_login"`
	Status            string     `json:"status" cql:"status"`
	CreatedAt         time.Time  `json:"created_at" cql:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" cql:"updated_at"`
}

func (g GithubInstallation) GetTable() *table.Table { return githubInstallationTable }
func (g GithubInstallation) PreCreate() error       { return nil }
func (g GithubInstallation) PreUpdate() error       { return nil }
