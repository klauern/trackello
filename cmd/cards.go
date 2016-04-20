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
	"github.com/spf13/cobra"
	"os"
	"log"
	"fmt"

	"github.com/VojtechVitek/go-trello"
)

// cardsCmd represents the cards command
var cardsCmd = &cobra.Command{
	Use:   "cards",
	Short: "List all of the cards on a particular board.",
	Long: `List all of the cards on a board.`,
	Run: listCardsOnBoard,
}

func init() {
	RootCmd.AddCommand(cardsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cardsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cardsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}


func listCardsOnBoard(cmd *cobra.Command, args []string) {
	conn, err := newTrelloConnection()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	
	if lists, err := conn.board.Lists(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		fmt.Println("Lists: ")
		for _, list := range lists {
			fmt.Printf("* %s\n", list.Name)
			printCards(list)
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