package activities

import (
	"context"
	"log/slog"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Notify sends chat notifications.
	Notify struct{}
)

// LinesExceeded notifies a chat service of exceeded lines. It uses the context and event to dispatch a
// notification via a chat hook. Returns error if notification fails, logging a warning.
func (n *Notify) LinesExceeded(ctx context.Context, evt *events.Event[eventsv1.ChatHook, eventsv1.Diff]) error {
	if err := kernel.Get().ChatHook(evt.Context.Hook).NotifyLinesExceed(ctx, evt); err != nil {
		slog.Warn("unable to notify on chat", "error", err.Error())
		return err
	}

	return nil
}

// MergeConflict notifies a chat service of a merge conflict. It uses the context and event to dispatch a
// notification via a chat hook. Returns error if notification fails, logging a warning.
func (n *Notify) MergeConflict(ctx context.Context, evt *events.Event[eventsv1.ChatHook, eventsv1.Merge]) error {
	if err := kernel.Get().ChatHook(evt.Context.Hook).NotifyMergeConflict(ctx, evt); err != nil {
		slog.Warn("unable to notify on chat", "error", err.Error())
		return err
	}

	return nil
}
