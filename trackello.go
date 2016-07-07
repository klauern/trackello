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
	stats         map[string]cardStatistics
}

type boardActions struct {
}

// cardStatistics provides a way to show various pieces of
type cardStatistics struct {
	comments, // represented by a horizontal ellepsis ⋯ 0x22EF
	updates, // represented by a keyboard 0x2328
	checklistsCreated, // represented by plus +
	checklistItemsChecked int // represented by check mark ✓ 0x2713
}

// Trackello represents the connection to Trello for a specific user.
type Trackello struct {
	token  string
	appKey string
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
	board, err := createTrelloBoardConnection(id)
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
	listActions := make(map[string][]string)
	for k, v := range activities.boardActions {
		 for _, vv := range v {
			 if listActions[vv.Data.List.Name] == nil {
				 listActions[vv.Data.List.Name] = []string{k}
			 } else {
				 listActions[vv.Data.List.Name] = append(listActions[vv.Data.List.Name], k)
			 }
		 }
	}

	for k, v := range listActions {
		fmt.Printf("-- %s\n", k)
		for _, vv := range v {
			fmt.Printf("   * %s\n", vv)
		}
	}
	//for k, v := range activities.boardActions {
	//	fmt.Printf("* %s\n", k)
	//	for _, vv := range v {
	//		fmt.Printf("  - %-24s %ss\n", vv.Date, vv.Type)
	//	}
	//}
}

func createTrelloBoardConnection(id string) (trello.Board, error) {
	conn, err := NewTrelloConnection()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var board trello.Board
	if id == "" {
		board, err = conn.PrimaryBoard()
	} else {
		board, err = conn.BoardWithId(id)
	}

	return board, err
}

func ListBoardActions(id string) error {
	board, err := createTrelloBoardConnection(id)
	if err != nil {
		return err
	}

	lists, err := board.Lists()
	if err != nil {
		return err
	}
	listMap := make(map[string]string)
	for _, v := range lists {
		listMap[v.Id] = v.Name
	}

	//actions, err := board.Actions(rest.CreateArgsForBoardActions()[0])
	//if err != nil {
	//	return err
	//}
	//for _, v := range actions {
	//	v.Data.Card.Id
	//	v.Data.Card.Name
	//
	//}

	//cards, err := board.Cards()
	//if err != nil {
	//	return err
	//}
	//for _, v := range cards {
	//}

	return nil
}
