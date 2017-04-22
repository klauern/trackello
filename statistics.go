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

package trackello

import (
	"fmt"

	"github.com/VojtechVitek/go-trello"
	"github.com/fatih/color"
)

// Statistics provides a way to show statistical information about a list, card or whatnot by aggregating the updates,
// comments, checklists, and other actions under a specific item.
type Statistics struct {
	comments, // represented by a horizontal ellepsis ⋯ 0x22EF
	updates, // represented by a keyboard 0x2328
	creates, // represented by plus +
	checkListItemUpdates int // represented by check mark ✓ 0x2713
}

// AddCalculation will the new action to the Card's statistics.
func (c *Card) AddCalculation(a trello.Action) error {
	switch a.Type {
	case "commentCard":
		c.stats.comments++
	case "updateCard", "deleteAttachmentFromCard", "updateList":
		c.stats.updates++
	case "createCard", "addChecklistToCard", "addAttachmentToCard":
		c.stats.creates++
	case "updateCheckItemStateOnCard":
		c.stats.checkListItemUpdates++
	default:
		c.stats.updates++
		return fmt.Errorf("unmapped action type: %s.  Defaulting to update", a.Type)
	}
	return nil
}

// PrintStatistics will print the statistics information out.
// Example format: [ 3 +  2 ≡  0 ✓  1 … ]
func (s *Statistics) PrintStatistics() string {
	if s == nil {
		s = &Statistics{}
	}
	stats := "[" + color.New(color.FgHiCyan).SprintfFunc()("%-2d +", s.updates)
	stats = stats + color.New(color.FgHiRed).SprintfFunc()(" %-2d ≡", s.comments)
	stats = stats + color.New(color.FgHiGreen).SprintfFunc()(" %-2d ✓", s.checkListItemUpdates)
	stats = stats + color.New(color.FgHiMagenta).SprintfFunc()(" %-2d …", s.creates)
	stats = stats + "]"
	return stats
}

// Total will print out the total number of actions that accumulate to make up the statistics information for a Card.
func (s *Statistics) Total() int {
	return s.checkListItemUpdates + s.comments + s.creates + s.updates
}
