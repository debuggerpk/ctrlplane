// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.


package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/slack-go/slack"
	"go.temporal.io/sdk/activity"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/providers/github"
)

type (
	// Activities groups all the activities for the slack provider.
	Activities struct{}
)

var (
	coreacts   *code.Activities
	githubacts *github.Activities
)

func (a *Activities) SendStaleBranchMessage(ctx context.Context, payload *defs.MessageIOStaleBranchPayload) error {
	logger := activity.GetLogger(ctx)

	token, err := reveal(payload.MessageIOPayload.BotToken, payload.MessageIOPayload.WorkspaceID)
	if err != nil {
		logger.Error("Error in reveal", "Error", err)
		return err
	}

	client, err := instance.GetSlackClient(token)
	if err != nil {
		logger.Error("Error in GetSlackClient", "Error", err)
		return err
	}

	attachment := formatStaleBranchAttachment(payload)

	// call blockset to send the message to slack channel or sepecific workspace.
	if err := notify(client, payload.MessageIOPayload.ChannelID, attachment); err != nil {
		logger.Error("Failed to post message to channel", "Error", err)
		return err
	}

	logger.Info("Slack notification sent successfully")

	return nil
}

func (a *Activities) SendNumberOfLinesExceedMessage(ctx context.Context, payload *defs.MessageIOLineExeededPayload) error {
	logger := activity.GetLogger(ctx)

	token, err := reveal(payload.MessageIOPayload.BotToken, payload.MessageIOPayload.WorkspaceID)
	if err != nil {
		logger.Error("Error in reveal", "Error", err)
		return err
	}

	client, err := instance.GetSlackClient(token)
	if err != nil {
		logger.Error("Error in GetSlackClient", "Error", err)
		return err
	}

	attachment := formatLineThresholdExceededAttachment(payload)

	// Call function to send the message to Slack channel or specific workspace.
	if err := notify(client, payload.MessageIOPayload.ChannelID, attachment); err != nil {
		logger.Error("Failed to post message to channel", "Error", err)
		return err
	}

	logger.Info("Slack notification sent successfully")

	return nil
}

func (a *Activities) SendMergeConflictsMessage(ctx context.Context, payload *defs.MergeConflictMessage) error {
	logger := activity.GetLogger(ctx)

	token, err := reveal(payload.MessageIOPayload.BotToken, payload.MessageIOPayload.WorkspaceID)
	if err != nil {
		logger.Error("Error in reveal", "Error", err)
		return err
	}

	client, err := instance.GetSlackClient(token)
	if err != nil {
		logger.Error("Error in GetSlackClient", "Error", err)
		return err
	}

	attachment := formatMergeConflictAttachment(payload)

	// call blockset to send the message to slack channel or sepecific workspace.
	if err := notify(client, payload.MessageIOPayload.ChannelID, attachment); err != nil {
		logger.Error("Failed to post message to channel", "Error", err)
		return err
	}

	logger.Info("Slack notification sent successfully")

	return nil
}

// TODO - move the uint functions.
func (a *Activities) NotifyMergeConflict(ctx context.Context, event *defs.Event[defs.MergeConflict, defs.RepoProvider]) error {
	logger := activity.GetLogger(ctx)

	token := ""
	target := ""

	var err error

	fields := a.conflict_fields(event)

	// by default we send the message to the channel, but if the user is a member of the team, we send the message to the user.
	// if we are sending on the channel, we introduce the owner of the branch to the channel.
	if event.Subject.UserID.String() != db.NullUUID {
		token, target, err = a.user_tokens(ctx, event)
		if err != nil {
			logger.Error("Error in token or target", "Error", err)
			return err
		}
	} else {
		token, target, err = a.repo_tokens(ctx, event)
		if err != nil {
			return err
		}

		owner := event.Payload.BaseCommit.Author
		url := fmt.Sprintf("https://github.com/%s", owner)
		fields = append(fields, slack.AttachmentField{
			Title: "Branch Owner",
			Value: fmt.Sprintf("[%s](%s)", owner, url),
			Short: false,
		})
	}

	attachment := slack.Attachment{
		Color: "warning",
		Pretext: fmt.Sprintf(`We've detected a merge conflict in your feature branch, <%s|%s>. 
	This means there are changes in your branch that clash with recent updates on the main branch (trunk).`,
			event.Payload.HeadCommit.URL, event.Payload.HeadCommit.SHA[:7]), // TODO
		Fallback:   "Merge Conflict Detected",
		MarkdownIn: []string{"fields"},
		Footer:     footer,
		Fields:     fields,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}

	client, err := instance.GetSlackClient(token)
	if err != nil {
		logger.Error("Error in GetSlackClient", "Error", err)
		return err
	}

	return notify(client, target, attachment)
}

func (a *Activities) conflict_fields(event *defs.Event[defs.MergeConflict, defs.RepoProvider]) []slack.AttachmentField {
	fields := []slack.AttachmentField{
		{
			Title: "*Branch*",
			Value: fmt.Sprintf("%s", event.Payload.BaseBranch),
			Short: true,
		}, {
			Title: "Current HEAD",
			Value: fmt.Sprintf("%s", event.Payload.HeadBranch),
			Short: true,
		}, {
			Title: "Conflict HEAD",
			Value: fmt.Sprintf("%s", event.Payload.HeadCommit.SHA),
			Short: true,
		},
		{
			Title: "Affected Files",
			Value: fmt.Sprintf("%s", FormatFilesList(event.Payload.Files)),
			Short: false,
		},
	}

	return fields
}

func (a *Activities) user_tokens(ctx context.Context, event *defs.Event[defs.MergeConflict, defs.RepoProvider]) (string, string, error) {
	tuser, err := auth.TeamUserIO().Get(ctx, event.Subject.UserID.String(), event.Subject.TeamID.String())

	if err != nil {
		return "", "", err
	}

	token, err := reveal(tuser.MessageProviderUserInfo.Slack.BotToken, tuser.MessageProviderUserInfo.Slack.ProviderTeamID)
	if err != nil {
		return "", "", err
	}

	return token, tuser.MessageProviderUserInfo.Slack.ProviderUserID, nil
}

func (a *Activities) repo_tokens(ctx context.Context, event *defs.Event[defs.MergeConflict, defs.RepoProvider]) (string, string, error) {
	repo, err := code.RepoIO().GetByID(ctx, event.Subject.ID.String())
	if err != nil {
		return "", "", err
	}

	token, err := reveal(repo.MessageProviderData.Slack.BotToken, repo.MessageProviderData.Slack.WorkspaceID)
	if err != nil {
		return "", "", err
	}

	return token, repo.MessageProviderData.Slack.ChannelID, nil
}
