package trackello

import (
	"fmt"

	"github.com/VojtechVitek/go-trello"
)

// Card is both the Trello Card + other stats on the actions in it.
type Card struct {
	card  *trello.Card
	stats *statistics
}

func (c *Card) String() string {
	return fmt.Sprintf("%s %s", c.stats.PrintStatistics(), c.card.Name)
}
