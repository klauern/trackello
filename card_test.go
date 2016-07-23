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
	if !sort.IsSorted(ByStatisticsCountRev(cards)) {
		t.Fatal("Not sorted")
	}
	sorted := map[int]int{
		0: 10,
		1: 5,
		2: 4,
		3: 0,
	}
	for i, v := range cards {
		if v.stats.Total() != sorted[i] {
			t.Fatalf("Error testing Sort: Expected %d, got %d", sorted[i], v.stats.Total())
		}
	}
}
