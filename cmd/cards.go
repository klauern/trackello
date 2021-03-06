// Copyright © 2018 Nick Klauer <klauer@gmail.com>
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

	"github.com/spf13/cobra"

	"github.com/VojtechVitek/go-trello"
	"github.com/klauern/trackello"
	"github.com/spf13/viper"
)

// cardsCmd represents the cards command
var cardsCmd = &cobra.Command{
	Use:   "cards",
	Short: "List all of the cards on a particular board.",
	Long:  `List all of the cards on a board.`,
	Run:   listCardsOnBoard,
}

func init() {
	RootCmd.AddCommand(cardsCmd)
}

func listCardsOnBoard(cmd *cobra.Command, args []string) {
	token := viper.GetString("token")
	appKey := viper.GetString("appkey")

	conn, err := trackello.NewTrackello(token, appKey)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	boardID := viper.GetString("board")
	if len(args) > 0 {
		boardID = args[0]
	}

	if board, err := conn.Board(boardID); err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		if lists, err := board.Lists(); err == nil {
			fmt.Println("Lists: ")
			for _, list := range lists {
				fmt.Printf("* %s\n", list.Name)
				printCards(list)
			}
		} else {
			log.Fatal(err)
			os.Exit(1)
		}
	}
}

func printCards(list trello.List) {
	cards, err := list.Cards()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	for _, card := range cards {
		fmt.Printf("  - %s\n", card.Name)
	}
}
