// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

	gotrello "github.com/VojtechVitek/go-trello"
	"github.com/klauern/trackello/trello"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var boardActions map[string][]gotrello.Action

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Track()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func Track() {
	token := viper.GetString(trello.API_TOKEN)
	appKey := viper.GetString(trello.API_APPKEY)

	// New Trello Client
	tr, err := gotrello.NewAuthClient(appKey, &token)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	board, err := tr.Board(viper.GetString(trello.PREFERRED_BOARD))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	args := trello.CreateArgsForBoardActions()
	actions, err := board.Actions(args...)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var cardsWorkedOn map[string]time.Time = make(map[string]time.Time)
	var oldestDate time.Time = time.Now()
	boardActions = make(map[string][]gotrello.Action)

	for _, action := range actions {
		switch boardActions[action.Data.Card.Name] {
		case nil:
			boardActions[action.Data.Card.Name] = []gotrello.Action{action}
		default:
			boardActions[action.Data.Card.Name] = append(boardActions[action.Data.Card.Name], action)
		}
		actionDate, err := time.Parse(trello.DateLayout, action.Date)
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
