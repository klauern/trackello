// Copyright © 2017 Nick Klauer <klauer@gmail.com>
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

	"github.com/klauern/trackello"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// boardsCmd represents the boards command
var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "List all of your boards",
	Long: `If you want to configure a default board, you will likely want to find out the
board ID for every one of them so you can specify this
somewhere (such as in your ~/.trackello.yaml).  This is pretty
easy.  Simply call 'boards' to get a listing of all the boards (real ID's omitted):

$ trackello boards
Using config file: ~/.trackello.yaml
Printing all OPEN Boards
Board Name                              ID
==========                              ==
FunThings                               XXXXXXXXXXXXXXXXXXXXXXXX
Family                                  XXXXXXXXXXXXXXXXXXXXXXXX
Personal Projects                       XXXXXXXXXXXXXXXXXXXXXXXX
Study                                   XXXXXXXXXXXXXXXXXXXXXXXX
Welcome Board                           XXXXXXXXXXXXXXXXXXXXXXXX
`,
	Run: listBoards,
}

func init() {
	RootCmd.AddCommand(boardsCmd)
}

func listBoards(cmd *cobra.Command, args []string) {
	token := viper.GetString("token")
	appKey := viper.GetString("appkey")

	trelloConn, err := trackello.NewTrackello(token, appKey)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println("Printing all OPEN Boards")
	boards, err := trelloConn.Boards()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("%-32s\t%-20s\n%-32s\t%-20s\n", "Board Name", "ID", "==========", "==")
	for _, v := range boards {
		if !v.Closed {
			fmt.Printf("%-32s\t%-30s\n", v.Name, v.Id)
		}
	}
}
