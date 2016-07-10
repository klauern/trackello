package trackello

import "github.com/VojtechVitek/go-trello"

// statistics provides a way to show statistical information about a list, card or whatnot by aggregating the updates,
// comments, checklists, and other actions under a specific item.
type statistics struct {
	comments, // represented by a horizontal ellepsis ⋯ 0x22EF
	updates, // represented by a keyboard 0x2328
	checklistsCreated, // represented by plus +
	checklistItemsChecked int // represented by check mark ✓ 0x2713
}

// Statistics represents the statistics for all the actions generated for a list, card, etc.
type Statistics interface {
	Calculate(trello.Action)
	PrintStatistics() string
}

func (c *Card) Calculate(a trello.Action) {}

func (c *Card) PrintStatistics() string {
	return ""
}

func (l *List) Calculate(a trello.Action) {}

func (l *List) PrintStatistics() string {
	return ""
}