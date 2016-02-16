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
	"os"

	"log"
	"time"

	"github.com/VojtechVitek/go-trello"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	dateLayout      string = "2006-01-02T15:04:05Z"
	API_APPKEY      string = "TRELLO_APPKEY"
	API_TOKEN       string = "TRELLO_TOKEN"
	PREFERRED_BOARD string = "preferredBoard"
)

var cfgFile string
var boardActions map[string][]trello.Action

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "trackello",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		Track()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.trackello.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}


	viper.SetConfigName(".trackello") // name of config file (without extension)
	if cwd, err := os.Getwd(); err == nil {
		viper.AddConfigPath(cwd)
	}
	viper.AddConfigPath("$HOME")      // adding home directory as first search path
	viper.AutomaticEnv()              // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}



func createArgsForBoardActions() []*trello.Argument {
	var args []*trello.Argument
	twoWeeksAgo := time.Now().Add(-1 * time.Hour * 24 * 14).Format(dateLayout)
	args = append(args, trello.NewArgument("since", twoWeeksAgo))
	args = append(args, trello.NewArgument("limit", "500"))
	return args
}

func Track() {
	token := viper.GetString(API_TOKEN)
	appKey := viper.GetString(API_APPKEY)

	// New Trello Client
	tr, err := trello.NewAuthClient(appKey, &token)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	board, err := tr.Board(viper.GetString(PREFERRED_BOARD))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	args := createArgsForBoardActions()
	actions, err := board.Actions(args...)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var cardsWorkedOn map[string]time.Time = make(map[string]time.Time)
	var oldestDate time.Time = time.Now()
	boardActions = make(map[string][]trello.Action)

	for _, action := range actions {
		switch boardActions[action.Data.Card.Name] {
		case nil:
			boardActions[action.Data.Card.Name] = []trello.Action{action}
		default:
			boardActions[action.Data.Card.Name] = append(boardActions[action.Data.Card.Name], action)
		}
		actionDate, err := time.Parse(dateLayout, action.Date)
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
