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
	if len(args) > 0 {
		trackello.Track(args[0])
	} else {
		trackello.Track(boardId)
	}
}
