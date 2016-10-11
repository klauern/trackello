package trackello

import (
	"reflect"
	"sort"
	"testing"

	"github.com/VojtechVitek/go-trello"
)

func newTestList() []List {
	return []List{
		*NewList(&trello.List{Name: "Z"}),
		*NewList(&trello.List{Name: "B"}),
		*NewList(&trello.List{Name: "A"}),
		*NewList(&trello.List{Name: "Q"}),
	}
}

func TestSortListByName(t *testing.T) {
	lists := newTestList()
	sort.Sort(ByListName(lists))
	if !sort.IsSorted(ByListName(lists)) {
		t.Fatal("Error: Not Sorted")
	}
	sortedLists := []List{
		*NewList(&trello.List{Name: "A"}),
		*NewList(&trello.List{Name: "B"}),
		*NewList(&trello.List{Name: "Q"}),
		*NewList(&trello.List{Name: "Z"}),
	}
	if !reflect.DeepEqual(lists, sortedLists) {
		t.Fatal("Lists are not equal")
	}
}

func TestPrintNonZeroActions(t *testing.T) {
	l := *NewList(&trello.List{Name: "test"})
	if l.PrintNonZeroActions() != "" {
		t.Fatalf("PrintNonZeroActions did not pass: expected %s, got %s", "", l.PrintNonZeroActions())
	}
}
