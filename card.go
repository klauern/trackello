// Copyright © 2016 Nick Klauer <klauer@gmail.com>
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

	"github.com/VojtechVitek/go-trello"
)

// Card is both the Trello Card + other stats on the actions in it.
type Card struct {
	card  *trello.Card
	stats *statistics
}

// this is a shortcut to mapping to a string since I may forget why it's a s
type cardID string

func (c *Card) String() string {
	return fmt.Sprintf("%s %s", c.stats.PrintStatistics(), c.card.Name)
}

// NewCard will construct a new Card element using a trello.Card as a base type.
func NewCard(c trello.Card) Card {
	return Card{
		card:  &c,
		stats: &statistics{},
	}
}

// ByStatisticsCountRev is a sortable type for the Card slice, sorting by the number of actions performed in descending order.
type ByStatisticsCountRev []Card

// Len returns the length of the underlying []Card slice.
func (c ByStatisticsCountRev) Len() int {
	return len(c)
}

// Swap swaps the positions of two Card items in the underlying []Card slice.
func (c ByStatisticsCountRev) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less returns whether the [i] position has MORE actions on it's underlying Card than the [j] element.
func (c ByStatisticsCountRev) Less(i, j int) bool {
	iTot := c[i].stats.updates + c[i].stats.checkListItemUpdates + c[i].stats.comments + c[i].stats.creates
	jTot := c[j].stats.updates + c[j].stats.checkListItemUpdates + c[j].stats.comments + c[j].stats.creates
	return iTot > jTot
}
