package trackello

import (
	"fmt"
	"log"
	"os"

	"github.com/VojtechVitek/go-trello"
	"github.com/klauern/trackello/rest"
	"github.com/spf13/viper"
)

// Trackello represents the connection to Trello for a specific user.
type Trackello struct {
	token  string
	appKey string
	client *trello.Client
}

// Card is both the Trello Card + other stats on the actions in it.
type Card struct {
	card  *trello.Card
	stats *statistics
}

// List is both the Trello List + other stats on the actions in it.
type List struct {
	name  string
	cards map[string]Card
	stats *statistics
}

// NewTrackello will create a `Trackello` type using your preferences application token and appkey.
func NewTrackello() (*Trackello, error) {
	token := viper.GetString("token")
	appKey := viper.GetString("appkey")

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

// PrimaryBoard will return a board based on your settings for a primary board.
func (t *Trackello) PrimaryBoard() (trello.Board, error) {
	board, err := t.client.Board(viper.GetString("board"))
	if err != nil {
		log.Fatal(err)
		return *board, err
	}
	return *board, nil
}

// BoardWithId will return the Trello Board given it's ID string.
func (t *Trackello) BoardWithId(id string) (trello.Board, error) {
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
				lc = newCardMap(card, listCards, list.Name)
				lc = listCards[list.Name]
			}
			if _, cok := lc.cards[card.Name]; !cok {
				newCard := Card{
					card: card,
				}
				lc.cards[card.Name] = newCard
			}
			c, _ := lc.cards[card.Name]
			c.AddCalculation(action)
			lc.cards[card.Name] = c
			listCards[list.Name] = lc
		}
	}
	return makeList(listCards), nil
}

func newCardMap(card trello.Card, listCards map[string]List, listName string) map[string]Card {
	cards := make(map[string]Card)
	cards[card.Name] = Card{card: card}
	listCards[listName] = List{
		name:  listName,
		cards: cards,
	}
	return cards
}

func makeList(listMap map[string]List) []List {
	list := make([]List, len(listMap))
	for _, v := range listMap {
		list = append(list, v)
	}
	return list
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

func (c *Card) String() string {
	return fmt.Sprintf("%s %s", c.card.Name, c.stats.PrintStatistics())
}

// Board will pull the Trello Board with an ID.  If id is "", it will pull it from the PrimaryBoard configuration setting.
func (t *Trackello) Board(id string) (trello.Board, error) {
	if id == "" {
		return t.PrimaryBoard()
	}
	return t.BoardWithId(id)
}

// BoardActions will retrieve a slice of trello.Action based on the Board ID.
func BoardActions(id string) ([]trello.Action, error) {
	t, err := NewTrackello()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

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
