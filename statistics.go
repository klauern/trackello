package trackello

import (
	"fmt"

	"github.com/VojtechVitek/go-trello"
	"github.com/fatih/color"
)

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
	AddCalculation(trello.Action)
}

func (c *Card) AddCalculation(a trello.Action) {
	switch a.Type {
	default:
		fmt.Printf("Type: %s\n", a.Type)
		c.stats.updates++
	}
}

// PrintStatistics will print the statistics information out.
// Example format: [ 3 ⋯  2 +  0 ✓  1 … ]
func (s *statistics) PrintStatistics() string {
	stats := "[" + color.CyanString(" %i ⋯", s.updates)
	stats = stats + color.RedString(" %i +", s.comments)
	stats = stats + color.GreenString(" %i ✓", s.checklistItemsChecked)
	stats = stats + color.MagentaString(" %i …", s.checklistsCreated)
	stats = stats + " ]"
	return stats
}

func (l *List) AddCalculation(a trello.Action) {}
