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

// trelloConnection repesents the connection to Trello and your preferred Board.
type trelloConnection struct {
	token  string
	appKey string
	board  trello.Board
}

// Track pulls all the latest activity from your Trello board given you've set the token, appkey, and preferred board
// ID to use.
func Track() {
	conn, err := newTrelloConnection()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	args := rest.CreateArgsForBoardActions()
	actions, err := conn.board.Actions(args...)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	allActivity := newTrelloActivity()

	mapActionsAndDates(actions, allActivity)

	printBoardActions(actions, allActivity)
}

func newTrelloConnection() (*trelloConnection, error) {
	token := viper.GetString("token")
	appKey := viper.GetString("appkey")
	// New Trello Client
	tr, err := trello.NewAuthClient(appKey, &token)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	board, err := tr.Board(viper.GetString("board"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &trelloConnection{
		token:  token,
		appKey: appKey,
		board:  *board,
	}, nil
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
	fmt.Printf("Cards Worked from %s to now:\n", activities.oldestDate.Format(time.ANSIC))
	for k, v := range activities.boardActions {
		fmt.Printf("* %s\n", k)
		for _, vv := range v {
			fmt.Printf("  - %-24s %ss\n", vv.Date, vv.Type)
		}
	}
}
