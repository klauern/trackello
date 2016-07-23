package trackello

import (
	"github.com/VojtechVitek/go-trello"
	"sort"
	"testing"
)

func TestSortListByName(t *testing.T) {
	lists := []List{
		*NewList(&trello.List{Name: "Z"}),
		*NewList(&trello.List{Name: "B"}),
		*NewList(&trello.List{Name: "A"}),
		*NewList(&trello.List{Name: "Q"}),
	}
	sort.Sort(ByListName(lists))
	if !sort.IsSorted(ByListName(lists)) {
		t.Fatal("Error: Not Sorted")
	}

	if lists[0].name != "A" {
		t.Fatal("Error: Not Sorted")
	}
	if lists[1].name != "B" {
		t.Fatal("Error: Not Sorted")
	}
	if lists[2].name != "Q" {
		t.Fatal("Error: Not Sorted")
	}
	if lists[3].name != "Z" {
		t.Fatal("Error: Not Sorted")
	}
}
