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
