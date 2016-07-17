package trackello

import (
	"sync"

	"github.com/VojtechVitek/go-trello"
	"github.com/klauern/trackello/rest"
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
		id:      b.Id,
		board:   b,
		listMux: &sync.RWMutex{},
		lists:   make(map[string]*List),
		//cardIdMux:      &sync.RWMutex{},
		//cardIdToListId: make(map[string]string),
	}
}

func (b *Board) Populate() error {
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
			if err = trackList.MapCards(); err != nil {
				return
			}
			// 2. return the list actions to return to the board
			b.listMux.Lock()
			b.lists[trackList.list.Id] = trackList
			b.listMux.Unlock()
		}(list)
	}
	wg.Wait()

	return nil
}

func (b *Board) MapActions() error {
	actions, err := b.board.Actions(rest.CreateArgsForBoardActions()...)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	for _, action := range actions {
		action := action
		wg.Add(1)
		go func(a trello.Action) {
			defer wg.Done()

			b.listMux.RLock()
			//b.lists[]
		}(action)
	}
	return nil
}

// Print will print the board actions out
func (b *Board) PrintActions() {
	for _, list := range b.lists {
		b.listMux.RLock()
		list.Print()
		b.listMux.RUnlock()
	}
}
