package trackello

import (
	"github.com/VojtechVitek/go-trello"
	"sort"
	"testing"
)

func TestSorting(t *testing.T) {
	cards := []Card{
		NewCard(trello.Card{}),
		NewCard(trello.Card{}),
		NewCard(trello.Card{}),
		NewCard(trello.Card{}),
	}
	cards[0].stats.comments = 5
	cards[3].stats.creates = 10
	cards[1].stats.checkListItemUpdates = 4
	sort.Sort(ByStatisticsCountRev(cards))
	sorted := map[int]int{
		0: 0,
		1: 4,
		2: 5,
		3: 10,
	}
	for i, v := range cards {
		if v.stats.Total() != sorted[i] {
			t.Fatal("Not sorted")
		}
	}
}
