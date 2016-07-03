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
	"github.com/klauern/trackello"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// boardsCmd represents the boards command
var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "List all of your boards",
	Long: `If you want to configure a default board, you will likely want to find out the
board ID for every one of them so you can specify this
somewhere (such as in your ~/.trackello.yaml).  This is pretty
easy.  Simply call 'boards' to get a listing of all the boards:

$ trackello boards
Your Boards
===========



`,
	Run: listBoards,
}

func init() {
	RootCmd.AddCommand(boardsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// boardsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// boardsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func listBoards(cmd *cobra.Command, args []string) {
	trelloConn, err := trackello.NewTrelloConnection()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("Printing all OPEN Boards\n")
	boards, err := trelloConn.Boards()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("%-32s\t%-20s\n%-32s\t%-20s\n", "Board Name", "ID", "==========", "==")
	for _, v := range boards {
		if !v.Closed {
			fmt.Printf("%-32s\t%-50s\n", v.Name, v.Id)
		}
	}
}
