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

package cmd

import (
	"fmt"

	"log"
	"os"
	"time"

	"github.com/VojtechVitek/go-trello"
	"github.com/klauern/trackello/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List activities on a board",
	Long: `List will pull all the activities for a particular
Trello board and list them in descending order.  This is useful
if you find yourself having to see what you've been working on`,
	Run: func(cmd *cobra.Command, args []string) {
		Track()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

// trelloConnection repesents the connection to Trello and your preferred Board.
type trelloConnection struct {
	token string
	appKey string
	board trello.Board
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
		token:token,
		appKey:appKey,
		board:*board,
	}, nil
}

// Track pulls all the latest activity from your Trello board given you've set the token, appkey, and preferred board
// ID to use.
// TODO: cmd\list.go:78::warning: cyclomatic complexity 12 of function Track() is high (> 10) (gocyclo)
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

	cardsWorkedOn := make(map[string]time.Time)
	oldestDate := time.Now()
	boardActions := make(map[string][]trello.Action)

	for _, action := range actions {
		switch boardActions[action.Data.Card.Name] {
		case nil:
			boardActions[action.Data.Card.Name] = []trello.Action{action}
		default:
			boardActions[action.Data.Card.Name] = append(boardActions[action.Data.Card.Name], action)
		}
		actionDate, err := time.Parse(rest.DateLayout, action.Date)
		if err != nil {
			continue // skip this one
		}
		if actionDate.Before(oldestDate) {
			oldestDate = actionDate
		}
		cardDate := cardsWorkedOn[action.Data.Card.Name]
		if cardDate.IsZero() || cardDate.After(actionDate) {
			cardsWorkedOn[action.Data.Card.Name] = actionDate
		}
	}

	fmt.Printf("Cards Worked from %s to now:\n", oldestDate.Format(time.ANSIC))
	for k, v := range boardActions {
		fmt.Printf("* %s\n", k)
		for _, vv := range v {
			fmt.Printf("  - %-24s %ss\n", vv.Date, vv.Type)
		}
	}
}
