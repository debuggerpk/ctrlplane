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

package core

import (
	"encoding/json"
	"errors"
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
)

func (repo *Repo) PreCreate() error { return nil }
func (repo *Repo) PreUpdate() error { return nil }

// TODO: move these entities to be generated by the code generator

var (
	changesetColums = []string{
		"id",
		"stack_id",
		"repo_markers",
		"created_by",
		"created_at",
		"updated_at",
	}

	changesetMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "changesets",
			Columns: changesetColums,
		},
	}

	changesetTable = itable.New(*changesetMeta.M)

	rolloutColumns = []string{
		"id",
		"stack_id",
		"blueprint_id",
		"trigger",
		"state",
		"created_at",
		"updated_at",
	}

	rolloutMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "rollouts",
			Columns: rolloutColumns,
		},
	}

	rolloutTable = itable.New(*rolloutMeta.M)
)

type (
	// ChangeSet records the state of the stack at a given point in time.
	// For a poly-repo BluePrint, a PR on one repo can trigger a release for the BluePrint.
	ChangeSet struct {
		ID          gocql.UUID           `json:"id" cql:"id"`
		StackID     gocql.UUID           `json:"stack_id" cql:"stack_id"`
		RepoMarkers ChangeSetRepoMarkers `json:"repo_markers" cql:"repo_markers"`
		CreatedBy   string               `json:"created_by" cql:"created_by"`
		CreatedAt   time.Time            `json:"created_at"`
		UpdatedAt   time.Time            `json:"updated_at"`
	}

	Rollout struct {
		ID          gocql.UUID   `json:"id" cql:"id"`
		StackID     gocql.UUID   `json:"stack_id" cql:"stack_id"`
		BlueprintID gocql.UUID   `json:"blueprint_id" cql:"blueprint_id"`
		ChangeSetID gocql.UUID   `json:"changeset_id" cql:"changeset_id"`
		State       RolloutState `json:"state" cql:"state"` // "in_progress" | "live" | "rejected"
		CreatedAt   time.Time    `json:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at"`
	}
)

func (changeset *ChangeSet) GetTable() itable.ITable { return changesetTable }
func (changeset *ChangeSet) PreCreate() error        { return nil }
func (changeset *ChangeSet) PreUpdate() error        { return nil }

func (rollout *Rollout) GetTable() itable.ITable { return rolloutTable }
func (rollout *Rollout) PreCreate() error        { return nil }
func (rollout *Rollout) PreUpdate() error        { return nil }

type (

	// RolloutState is the state of a rollout.
	RolloutState        string
	RolloutStateMapType map[string]RolloutState

	ChangeSetRepoMarker struct {
		Provider   string `json:"provider"`
		CommitID   string `json:"commit_id"`
		RepoID     string `json:"repo_id"`
		HasChanged bool   `json:"changed"`
	}

	ChangeSetRepoMarkers []ChangeSetRepoMarker
)

const (
	RolloutStateQueued     RolloutState = "queued"
	RolloutStateInProgress RolloutState = "in_progress"
	RolloutStateCompleted  RolloutState = "completed"
	RolloutStateRejected   RolloutState = "rejected"
)

var (
	RolloutStateMap = RolloutStateMapType{
		RolloutStateQueued.String():     RolloutStateQueued,
		RolloutStateInProgress.String(): RolloutStateInProgress,
		RolloutStateCompleted.String():  RolloutStateCompleted,
		RolloutStateRejected.String():   RolloutStateRejected,
	}
)

func (marker ChangeSetRepoMarkers) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(marker)
}

func (marker *ChangeSetRepoMarkers) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, marker)
}

func (mp MessageProviderData) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(mp)
}

func (mp *MessageProviderData) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	if len(data) == 0 {
		*mp = MessageProviderData{}
		return nil
	}

	return json.Unmarshal(data, mp)
}

func (rs RolloutState) String() string {
	return string(rs)
}

func (rs RolloutState) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.String())
}

func (rs *RolloutState) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	val, ok := RolloutStateMap[s]
	if !ok {
		return errors.New("invalid rollout state")
	}

	*rs = val

	return nil
}

func (rs RolloutState) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(rs)
}

func (rs *RolloutState) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, rs)
}

func (rp RepoProvider) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return gocql.Marshal(info, rp.String())
}

func (rp *RepoProvider) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	*rp = RepoProviderMap[string(data)]

	return nil
}
