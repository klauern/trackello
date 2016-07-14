package trackello

import "fmt"

// List is both the Trello List + other stats on the actions in it.
type List struct {
	name  string
	cards map[string]Card
	stats *statistics
}

// Print will print out a list and all of the cards to the command-line.
func (l *List) Print() {
	if len(l.name) > 0 {
		fmt.Printf("%s\n", l.name)
		for _, card := range l.cards {
			fmt.Printf("  * %s\n", card.String())
		}
	}
}

func makeList(listMap map[string]List) []List {
	list := make([]List, len(listMap))
	for _, v := range listMap {
		list = append(list, v)
	}
	return list
}

// ByListName is a super type of a List[], with functions that implement the sort interface.
type ByListName []List

// Len returns the length of the ByListName slice.
func (l ByListName) Len() int {
	return len(l)
}

// Swap will swap the positions of two trackello.List items in the ByListName slice.
func (l ByListName) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Less determines which of the two trackello.List items is before other based on the List Name.
func (l ByListName) Less(i, j int) bool {
	return l[i].name < l[j].name
}
