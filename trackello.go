package trackello

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/VojtechVitek/go-trello"
	"github.com/klauern/trackello/rest"
	"github.com/spf13/viper"
)

type trelloActivity struct {
	cardsWorkedOn map[string]time.Time
	oldestDate    time.Time
	boardActions  map[string][]trello.Action
}

// Trackello represents the connection to Trello for a specific user.
type Trackello struct {
	token  string
	appKey string
	client trello.Client
}

// Card is both the Trello Card + other stats on the actions in it.
type Card struct {
	card  trello.Card
	stats statistics
}

// List is both the Trello List + other stats on the actions in it.
type List struct {
	cards []Card
	stats statistics
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
		client: *tr,
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
		return make([]trello.Board, 0), err
	}
	boards, err := member.Boards()
	return boards, err
}

// Track pulls all the latest activity from your Trello board given you've set the token, appkey, and preferred board
// ID to use.
func Track(id string) {
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

	allActivity := newTrelloActivity()
	t.mapActionsAndDates(actions, allActivity)
	fmt.Printf("Listing cards worked on \"%s\" for from %s to now:\n", board.Name, allActivity.oldestDate.Format(time.ANSIC))
	printBoardActions(actions, allActivity)
}

func newTrelloActivity() *trelloActivity {
	return &trelloActivity{
		cardsWorkedOn: make(map[string]time.Time),
		oldestDate:    time.Now(),
		boardActions:  make(map[string][]trello.Action),
	}
}

func (t *Trackello) mapActionsAndDates(actions []trello.Action, activities *trelloActivity) {
	for _, action := range actions {
		switch activities.boardActions[action.Data.Card.Name] {
		case nil:
			activities.boardActions[action.Data.Card.Name] = []trello.Action{action}
		default:
			activities.boardActions[action.Data.Card.Name] = append(activities.boardActions[action.Data.Card.Name], action)
		}
		if action.Data.List.Name == "" {
			action.Data.List.Name = t.getListForAction(action)
		}
	}
}

func (t *Trackello) getListForAction(a trello.Action) string {
	card, err := t.client.Card(a.Data.Card.Id)
	if err != nil {
		return ""
	}
	list, err := t.client.List(card.IdList)
	if err != nil {
		return ""
	}
	return list.Name
}

func printBoardActions(actions []trello.Action, activities *trelloActivity) {
	listActions := make(map[string][]string)
	for k, v := range activities.boardActions {
		for _, vv := range v {
			switch {
			case listActions[vv.Data.List.Name] == nil:
				listActions[vv.Data.List.Name] = []string{k}
			default:
				listActions[vv.Data.List.Name] = append(listActions[vv.Data.List.Name], k)
			}
		}
	}

	for k, v := range listActions {
		fmt.Printf("%s\n", k)
		for _, vv := range v {
			fmt.Printf("  * %s\n", vv)
		}

	}
}

// Board will pull the Trello Board with an ID.  If id is "", it will pull it from the PrimaryBoard configuration setting.
func (t *Trackello) Board(id string) (trello.Board, error) {
	if id == "" {
		return t.PrimaryBoard()
	}
	return t.BoardWithId(id)
}