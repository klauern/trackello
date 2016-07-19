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

	"github.com/klauern/trackello"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [board id]",
	Short: "List activities on a board",
	Long: `List will pull all the activities for a particular
Trello board and list them in descending order.  This is useful
if you find yourself having to see what you've been working on`,
	Run: ListActivity,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

// ListActivity will list all the card actions for a board, sorting by List.
func ListActivity(cmd *cobra.Command, args []string) {
	switch {
	case len(args) > 0:
		fmt.Printf("Printing Board activity for Board ID %s\n", args[0])
		PrintParallelBoardActivity(args[0])
	case len(viper.GetString("board")) > 0:
		fmt.Printf("Printing Board activity for Board ID %s\n", viper.GetString("board"))
		PrintParallelBoardActivity(viper.GetString("board"))
	default:
		panic("No board id specified in either .trackello.yaml or on command-line.")
	}
}

// PrintParallelBoardActivity will take a Board ID and print all of the activity that the board has, by parallelizing
// requests by List.
func PrintParallelBoardActivity(id string) {
	token := viper.GetString("token")
	appKey := viper.GetString("appkey")

	t, err := trackello.NewTrackello(token, appKey)
	if err != nil {
		panic(err)
	}
	b, err := t.Board(id)
	if err != nil {
		panic(err)
	}
	board := trackello.NewBoard(&b)
	err = board.PopulateLists()
	if err != nil {
		panic(err)
	}
	err = board.MapActions()
	if err != nil {
		panic(err)
	}

	board.PrintActions()
}