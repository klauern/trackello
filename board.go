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
	"sort"
	"sync"

	"fmt"
	"github.com/VojtechVitek/go-trello"
	"github.com/pkg/errors"
)

// Board is a super-type for a Trello board.  Board also contains a mutex and map of a List ID to a List.
type Board struct {
	id      string
	board   *trello.Board
	listMux *sync.RWMutex
	lists   map[string]List
}

// NewBoard will create a new Board type, using a trello.Board as a starting point.
func NewBoard(b *trello.Board) *Board {
	return &Board{
		id:      b.Id,
		board:   b,
		listMux: &sync.RWMutex{},
		lists:   make(map[string]List),
	}
}

// PopulateLists will Populate the board's lists with cards and missing data.
func (b *Board) PopulateLists() error {
	lists, err := b.board.Lists()
	if err != nil {
		return errors.Wrapf(err, "Unable to get Lists for Board %s", b.board.Name)
	}
	wg := sync.WaitGroup{}
	for _, list := range lists {
		list := list
		wg.Add(1)
		go func(list trello.List) {
			defer wg.Done()
			// 1. calculate the actions for the list
			trackList := NewList(&list)
			if err := trackList.MapCards(); err != nil {
				return
			}
			// 2. return the list actions to return to the board
			b.listMux.Lock()
			b.lists[trackList.list.Id] = *trackList
			b.listMux.Unlock()
		}(list)
	}
	wg.Wait()

	return nil
}

// MapActions queries Trello's API for all of the recent actions performed on a Board, and maps that to the
// board itself, into a list and card.
func (b *Board) MapActions() error {
	wg := sync.WaitGroup{}
	for _, list := range b.lists {
		wg.Add(1)
		go func(l List) {
			defer wg.Done()
			if err := l.MapActions(); err != nil {
				panic(err)
			}
		}(list)
	}
	wg.Wait()
	return nil
}

// PrintActions will print the board actions out.
func (b *Board) PrintActions() {
	lists := make([]List, len(b.lists))
	for _, list := range b.lists {
		lists = append(lists, list)
	}
	b.listMux.Lock()
	sort.Sort(ByListName(lists))
	b.listMux.Unlock()
	for _, list := range lists {
		b.listMux.RLock()
		fmt.Printf("%s", list.PrintNonZeroActions())
		b.listMux.RUnlock()
	}
}
