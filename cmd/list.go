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
	"github.com/klauern/trackello"
	"github.com/spf13/cobra"
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

var boardId string

func init() {
	RootCmd.AddCommand(listCmd)
}

func ListActivity(cmd *cobra.Command, args []string) {
	switch {
	case len(args) > 0:

		actions, err := trackello.BoardActions(args[0])
		if err != nil {
			panic(err)
		}
		t, err := trackello.NewTrackello()
		if err != nil {
			panic(err)
		}
		list, err := t.MapBoardActions(actions)
		if err != nil {
			panic(err)
		}

		for _, v := range list {
			v.Print()
		}
		//fmt.Printf("Listing cards worked on \"%s\" for from %s to now:\n", board.Name, allActivity.oldestDate.Format(time.ANSIC))

		//trackello.Track(args[0])
	case boardId != "":
		_, err := trackello.BoardActions(boardId)
		if err != nil {
			panic(err)
		}
		//trackello.Track(boardId)
	default:
		panic("No board id specified in either boardId or on command-line.")
	}

	// pseudocode for listing things
	//
	// 1. create connection
	// 2. get board
	// 4. get actions on board
	// -- in parallel
	//    #. map action to card
	//		 - add calculation to statistics
	//    #. map action to list
	//		 - add calculation to statistics
	// 6. map cards to lists
}
