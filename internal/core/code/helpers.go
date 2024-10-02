package code

import (
	"github.com/gocql/gocql"
)

type (
	// BranchTriggers maps branch names to their corresponding triggering event IDs.
	//
	// This data structure facilitates event lineage tracing by providing the root event for each branch.
	BranchTriggers map[string]gocql.UUID

	// StashedEvents stores events that are awaiting processing.
	//
	// Events are typically stashed when the associated branch does not yet exist or the event requires a
	// parent event (e.g., a push event needing a branch creation event) that has not yet been received. This
	// scenario can arise due to the distributed nature of event arrival.
	StashedEvents map[string][]RepoEvent
)

// add associates a branch with its triggering event ID.
func (b BranchTriggers) add(branch string, id gocql.UUID) {
	b[branch] = id
}

// del removes the association between a branch and its triggering event ID.
func (b BranchTriggers) del(branch string) {
	delete(b, branch)
}

// get retrieves the event ID associated with a branch.
//
// Returns the event ID and a boolean indicating whether the branch exists.
func (b BranchTriggers) get(branch string) (gocql.UUID, bool) {
	id, ok := b[branch]

	return id, ok
}

// push adds an event to the stash for the specified branch.
func (s StashedEvents) push(branch string, event RepoEvent) {
	if _, ok := s[branch]; !ok {
		s[branch] = make([]RepoEvent, 0)
	}

	s[branch] = append(s[branch], event)
}

// pop retrieves and removes the oldest event from the stash for the specified branch.
//
// Returns the event and a boolean indicating whether an event was present.
func (s StashedEvents) pop(branch string) (RepoEvent, bool) {
	events, ok := s[branch]
	if !ok || len(events) == 0 {
		return nil, false
	}

	event := events[0]
	s[branch] = events[1:]

	return event, true
}
