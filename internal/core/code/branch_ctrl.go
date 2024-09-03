package code

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
)

// BranchCtrl is the event loop to process events during the lifecycle of a branch.
//
// It processes the following events:
//
//   - push
//   - rebase
//   - create_delete
//   - pr
func BranchCtrl(ctx workflow.Context, repo *defs.Repo, branch string) error {
	selector := workflow.NewSelector(ctx)
	ctx, state := NewBranchCtrlState(ctx, repo, branch)

	// start the stale check coroutine.
	state.check_stale(ctx)

	// setup signals

	// push event signal.
	// detect changes. if changes are greater than threshold, send early warning message.
	push := workflow.GetSignalChannel(ctx, defs.RepoIOSignalPush.String())
	selector.AddReceive(push, state.on_push(ctx))

	// rebase signal.
	// attempts to rebase the branch with the base branch. if there are merge conflicts, sends message.
	rebase := workflow.GetSignalChannel(ctx, defs.RepoIOSignalRebase.String())
	selector.AddReceive(rebase, state.on_rebase(ctx))

	// create_delete signal.
	// creates or deletes the branch.
	create_delete := workflow.GetSignalChannel(ctx, defs.RepoIOSignalCreateOrDelete.String())
	selector.AddReceive(create_delete, state.on_create_delete(ctx))

	// pr signal.
	pr := workflow.GetSignalChannel(ctx, defs.RepoIOSignalPullRequestOpenedOrClosedOrReopened.String())
	selector.AddReceive(pr, state.on_pr(ctx))

	// label signal.
	lebal := workflow.GetSignalChannel(ctx, defs.RepoIOSignalPullRequestLabeledOrUnlabeled.String())
	selector.AddReceive(lebal, state.on_label(ctx))

	// main event loop
	for state.is_active() {
		selector.Select(ctx)

		// TODO - need to optimize
		// TODO - remove
		if state.pr != nil {
			_ctx, q_state := NewQueueCtrlState(ctx, repo, branch)
			q_state.push(_ctx, state.pr, false) // TODO - handle priority
		}

		if state.needs_reset() {
			return state.as_new(ctx, "event history exceeded threshold", BranchCtrl, repo, branch)
		}
	}

	// graceful shutdown
	state.terminate(ctx)

	return nil
}
