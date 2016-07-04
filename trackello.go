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
	//board  trello.Board
	client trello.Client
}

// NewTrelloConnection will create a `Trackello` type using your preferences application token and appkey.
func NewTrelloConnection() (*Trackello, error) {
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
	conn, err := NewTrelloConnection()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	args := rest.CreateArgsForBoardActions()
	var board trello.Board
	if id == "" {
		board, err = conn.PrimaryBoard()
	} else {
		board, err = conn.BoardWithId(id)
	}

	if err != nil {
		panic(err)
	}
	actions, err := board.Actions(args...)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	allActivity := newTrelloActivity()
	mapActionsAndDates(actions, allActivity)
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

func mapActionsAndDates(actions []trello.Action, activities *trelloActivity) {
	for _, action := range actions {
		switch activities.boardActions[action.Data.Card.Name] {
		case nil:
			activities.boardActions[action.Data.Card.Name] = []trello.Action{action}
		default:
			activities.boardActions[action.Data.Card.Name] = append(activities.boardActions[action.Data.Card.Name], action)
		}
		actionDate, err := time.Parse(rest.DateLayout, action.Date)
		if err != nil {
			continue // skip this one
		}
		if actionDate.Before(activities.oldestDate) {
			activities.oldestDate = actionDate
		}
		cardDate := activities.cardsWorkedOn[action.Data.Card.Name]
		if cardDate.IsZero() || cardDate.After(actionDate) {
			activities.cardsWorkedOn[action.Data.Card.Name] = actionDate
		}
	}
}

func printBoardActions(actions []trello.Action, activities *trelloActivity) {
	for k, v := range activities.boardActions {
		fmt.Printf("* %s\n", k)
		for _, vv := range v {
			fmt.Printf("  - %-24s %ss\n", vv.Date, vv.Type)
		}
	}
}
