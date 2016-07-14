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
	"log"
	"os"

	"github.com/VojtechVitek/go-trello"
	"github.com/klauern/trackello/rest"
)

// Trackello represents the connection to Trello for a specific user.
type Trackello struct {
	token  string
	appKey string
	client *trello.Client
}

// NewTrackello will create a `Trackello` type using your preferences application token and appkey.
func NewTrackello(token, appKey string) (*Trackello, error) {
	// New Trello Client
	tr, err := trello.NewAuthClient(appKey, &token)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Trackello{
		token:  token,
		appKey: appKey,
		client: tr,
	}, nil
}

// Board will return the Trello Board given it's ID string.
func (t *Trackello) Board(id string) (trello.Board, error) {
	board, err := t.client.Board(id)
	if err != nil {
		log.Fatal(err)
		return *board, err
	}
	return *board, nil
}

// Boards will list all of the boards for the authenticated user (i.e. 'me').
func (t *Trackello) Boards() ([]trello.Board, error) {
	member, err := t.client.Member("me")
	if err != nil {
		log.Fatalf("Error getting 'me' Member: %v", err)
		return make([]trello.Board, 0), err
	}
	boards, err := member.Boards()
	return boards, err
}

func (t *Trackello) getCardForAction(a trello.Action) (*trello.Card, error) {
	return t.client.Card(a.Data.Card.Id)
}

// MapBoardActions takes the slice of []trello.Action and maps it to a Card and it's associated List.
func (t *Trackello) MapBoardActions(actions []trello.Action) ([]List, error) {
	listCards := make(map[string]List)
	for _, action := range actions {
		if len(action.Data.Card.Id) > 0 {
			// TODO: restructure this into a map of mapped Cards
			// We can parallelize the calls, AND reduce making redundant calls.
			// trackello.List should have a map[cardID]Card instead of just a slice of []Card
			//
			card, err := t.getCardForAction(action)
			if err != nil {
				return nil, err
			}
			list, err := t.client.List(card.IdList)
			if err != nil {
				return nil, err
			}
			lc, ok := listCards[list.Name]
			if !ok {
				cards := make(map[cardId]Card)
				cards[cardId(card.Id)] = Card{
					card:  card,
					stats: &statistics{},
				}
				listCards[list.Name] = List{
					name:  list.Name,
					cards: cards,
				}
				lc = listCards[list.Name]
			}
			if _, cok := lc.cards[cardId(card.Id)]; !cok {
				newCard := Card{
					card:  card,
					stats: &statistics{},
				}
				lc.cards[cardId(card.Id)] = newCard
			}
			c, _ := lc.cards[cardId(card.Id)]
			c.AddCalculation(action)
			lc.cards[cardId(card.Id)] = c
			listCards[list.Name] = lc
		}
	}
	return makeList(listCards), nil
}

// BoardActions will retrieve a slice of trello.Action based on the Board ID.
func (t *Trackello) BoardActions(id string) ([]trello.Action, error) {
	board, err := t.Board(id)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	args := rest.CreateArgsForBoardActions()
	actions, err := board.Actions(args...)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return actions, err
}
