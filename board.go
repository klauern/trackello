package trackello

import (
	"sync"

	"github.com/VojtechVitek/go-trello"
	"github.com/pkg/errors"
)

type Board struct {
	id             string
	board          *trello.Board
	listMux        *sync.RWMutex
	lists          map[string]*List
	cardIdMux      *sync.RWMutex
	cardIdToListId map[string]string
}

func NewBoard(b *trello.Board) *Board {
	return &Board{
		id:             b.Id,
		board:          b,
		listMux:        &sync.RWMutex{},
		lists:          make(map[string]*List),
		cardIdMux:      &sync.RWMutex{},
		cardIdToListId: make(map[string]string),
	}
}

func (b *Board) Populate() error {
	lists, err := b.board.Lists()
	if err != nil {
		return errors.Wrapf(err, "Unable to get Lists for Board %s", b.board.Name)
	}
	wg := sync.WaitGroup{}
	for _, list := range lists {
		wg.Add(1)
		go func(list *trello.List) {
			defer wg.Done()
			// 1. calculate the actions for the list
			listActions := NewList(list)
			if err = listActions.MapCards(); err != nil {
				return
			}
			_, err := listActions.MapActions()
			if err != nil {
				return
			}
			// 2. return the list actions to return to the board
			b.listMux.Lock()
			b.lists[listActions.list.Id] = listActions
			b.listMux.Unlock()
		}(&list)
	}
	wg.Wait()
	return nil
}

// Print will print the board actions out
func (b *Board) PrintActions() {
	b.listMux.RLock()
	for _, list := range b.lists {
		list.Print()
	}
	b.listMux.RUnlock()
}
