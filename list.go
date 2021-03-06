// Copyright © 2018 Nick Klauer <klauer@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trackello

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/VojtechVitek/go-trello"
	"github.com/klauern/trackello/rest"
	"github.com/pkg/errors"
)

// List is both the Trello List + other stats on the actions in it.
type List struct {
	name  string
	cards map[cardID]Card
	stats *Statistics
	list  *trello.List
	mux   *sync.RWMutex
}

// PrintNonZeroActions will print out a list and all of the cards to the command-line that have
// more than 0 actions associated with them.
func (l *List) PrintNonZeroActions() string {
	cards := l.NonZeroActions()
	output := ""
	if len(cards) > 0 {
		output += fmt.Sprintf("%s\n", l.name)
		for _, card := range cards {
			output += fmt.Sprintf("  * %s\n", card.String())
		}
	}
	return output
}

// NonZeroActions returns a list of Card entries that have activity in them.
func (l *List) NonZeroActions() []Card {
	var cards []Card
	if len(l.name) > 0 {
		for _, card := range l.cards {
			if card.stats.Total() > 0 {
				cards = append(cards, card)
			}
		}
		sort.Sort(ByStatisticsCountRev(cards))
	}
	return cards
}

// NewList constructs a new *List based off of a go-trello *List.
func NewList(l *trello.List) *List {
	return &List{
		name:  l.Name,
		cards: make(map[cardID]Card),
		stats: &Statistics{},
		list:  l,
		mux:   &sync.RWMutex{},
	}
}

// MapActions will map all of the Actions that occurred on a List.
func (l *List) MapActions() error {
	// This is ugly.  I admit it.  It was a messy hack to get parallelism working for mapping
	// Trello card actions, and now I need to tease this apart, much like you would if you found
	// that a rat had been chewing through your pantry.  Separate out the stuff the rat's been in
	// with the good stuff that it didn't chew through. Ugh.
	args := rest.CreateArgsForBoardActions()
	actions, err := l.list.Actions(args...)
	if err != nil {
		return errors.Wrapf(err, "Error in MapActions getting List \"%s\" Actions: ", l.name)
	}
	for _, action := range actions {
		card, ok := l.cards[cardID(action.Data.Card.Id)]
		if !ok {
			switch action.Type {
			// Ignore list actions.  We're only interested in actions on Cards themselves.
			case "updateList", "createList":
				continue
			case "updateCard":
				// if we're moving cards between lists, we just won't map it.  It will likely be either
				// caught from the other list's actions, or not, but it's not worth digging too deeply
				if len(action.Data.ListBefore.Id) > 0 && len(action.Data.ListAfter.Id) > 0 {
					continue
				}
			}
		} else {
			if err := card.AddCalculation(action); err != nil {
				// If there's an error, it's probably because it's unmapped.  We may want to output that.
				fmt.Printf("%s\n", err)
			}
			l.mux.Lock()
			l.cards[cardID(action.Data.Card.Id)] = card
			l.mux.Unlock()
		}
	}
	return nil
}

// MapCards maps all of the cards for a list into the List.cards map[string]Card based on the cardID.
func (l *List) MapCards() error {
	cards, err := l.list.Cards()
	if err != nil {
		return errors.Wrapf(err, "Error in listing Cards in MapCards")
	}
	for _, card := range cards {
		l.cards[cardID(card.Id)] = NewCard(card)
	}
	return nil
}

// ByListName is a sortable type for []List, allowing sorting by the List Name (case-insensitive).
type ByListName []List

// Len returns the length of the ByListName slice.
func (l ByListName) Len() int {
	return len(l)
}

// Swap will swap the positions of two trackello.List items in the ByListName slice.
func (l ByListName) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Less determines which of the two trackello.List items is before other based on the List Name (lowercased).
func (l ByListName) Less(i, j int) bool {
	return strings.Compare(strings.ToLower(l[i].name), strings.ToLower(l[j].name)) == -1
}
