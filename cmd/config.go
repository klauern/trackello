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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
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
	viper.BindEnv("appkey", "TRACKELLO_TRELLO_APPKEY")
	viper.BindEnv("token", "TRACKELLO_TRELLO_TOKEN")
	viper.BindEnv("board", "TRACKELLO_TRELLO_PREFERREDBOARD")
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
	viper.BindPFlag("appkey", RootCmd.PersistentFlags().Lookup("appkey"))
	RootCmd.PersistentFlags().StringVar(&trelloToken, "token", "", "Trello Token")
	viper.BindPFlag("token", RootCmd.PersistentFlags().Lookup("token"))
	RootCmd.PersistentFlags().StringVar(&preferredBoard, "board", "", "Preferred Board ID")
	viper.BindPFlag("board", RootCmd.PersistentFlags().Lookup("board"))
	viper.RegisterAlias("preferredBoard", "board")
}

func runConfigCmd(cmd *cobra.Command, args []string) {
	// TODO: Work your own magic here
	fmt.Println("config called")

	fmt.Printf("Token is %s\n", viper.GetString("token"))
	fmt.Printf("AppKey is %s\n", viper.GetString("appkey"))
	fmt.Printf("Board ID is %s\n", viper.GetString("board"))

}
