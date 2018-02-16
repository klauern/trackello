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
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

// TrackelloVersion is the statically defined version of this project.
var TrackelloVersion string

var Commit string
var BuildDate string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Trackello's Version",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Trackello v%s\nBuildDate: %s\nCommit: %s\n", TrackelloVersion, BuildDate, Commit)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func setBuildDate() {
	fname, err := osext.Executable()
	if err != nil {
		fname = "trackello"
	}
	dir, err := filepath.Abs(filepath.Dir(fname))
	if err != nil {
		fmt.Println(err)
		return
	}
	fi, err := os.Lstat(filepath.Join(dir, filepath.Base(fname)))
	if err != nil {
		fmt.Println(err)
		return
	}
	t := fi.ModTime()
	BuildDate = t.Format(time.RFC3339)
}

func formatBuildDate() {
	t, err := time.Parse("2006-01-02T15:04:05-0700", BuildDate)
	if err != nil {
		t = time.Now()
	}
	BuildDate = t.Format(time.RFC3339)
}
