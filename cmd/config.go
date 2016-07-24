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
	"strings"

	"github.com/klauern/trackello"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	trelloAppKey,
	trelloToken,
	cfgFile,
	preferredBoard string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure your Trello API keys and preferred Trello Board",
	Long: `To communicate to a Trello API, you will need to configure a
minimum of 3 parameters:
  - TRELLO_APPKEY
  - TRELLO_TOKEN
  - preferredBoard

If you do not know or have any of these, you can review  the documentation
on this site: https://trello.com/app-key
`,
	Run: runConfigCmd,
}

// TODO: Need to look at whether this is just too much going on.
func init() {
	// enable ability to specify config file via flag
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.trackello.yaml)")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	viper.SetConfigName(".trackello") // name of config file (without extension)

	// Set Environment Variables
	viper.SetEnvPrefix("trackello")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_")) // replace environment variables to underscore (_) from hyphen (-)
	if err := viper.BindEnv("appkey", trackello.TRACKELLO_APPKEY); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("token", trackello.TRACKELLO_TOKEN); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("board", "TRACKELLO_TRELLO_PREFERREDBOARD"); err != nil {
		panic(err)
	}
	viper.AutomaticEnv() // read in environment variables that match every time Get() is called

	// Add Configuration Paths
	if cwd, err := os.Getwd(); err == nil {
		viper.AddConfigPath(cwd)
	}
	viper.AddConfigPath("$HOME") // adding home directory as first search path

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	RootCmd.AddCommand(configCmd)

	RootCmd.PersistentFlags().StringVar(&trelloAppKey, "appkey", "", "Trello Application Key")
	if err := viper.BindPFlag("appkey", RootCmd.PersistentFlags().Lookup("appkey")); err != nil {
		panic(err)
	}
	RootCmd.PersistentFlags().StringVar(&trelloToken, "token", "", "Trello Token")
	if err := viper.BindPFlag("token", RootCmd.PersistentFlags().Lookup("token")); err != nil {
		panic(err)
	}
	RootCmd.PersistentFlags().StringVar(&preferredBoard, "board", "", "Preferred Board ID")
	if err := viper.BindPFlag("board", RootCmd.PersistentFlags().Lookup("board")); err != nil {
		panic(err)
	}
	viper.RegisterAlias("preferredBoard", "board")
}

func runConfigCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("Token is %s\n", viper.GetString("token"))
	fmt.Printf("AppKey is %s\n", viper.GetString("appkey"))
	fmt.Printf("Board ID is %s\n", viper.GetString("board"))
}
