package trackello

import (
	"github.com/VojtechVitek/go-trello"
	"reflect"
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
