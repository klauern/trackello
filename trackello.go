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

package trackello

import (
	"fmt"
	"log"

	"github.com/VojtechVitek/go-trello"
)

const (
	// TrackelloAppKey represents the underlying environment variable for Trello's AppKey
	TrackelloAppKey = "TRACKELLO_TRELLO_APPKEY" // #nosec
	// TrackelloToken represents the underlying environment variables for the Trello Token
	TrackelloToken = "TRACKELLO_TRELLO_TOKEN" // #nosec
	// TrackelloPreferredBoard represents the underlying environment variable for the preferred Board ID to retrieve
	TrackelloPreferredBoard = "TRACKELLO_TRELLO_PREFERREDBOARD"
)

// Trackello represents the connection to Trello for a specific user.
type Trackello struct {
	token  string
	appKey string
	client *trello.Client
}

// NewTrackello will create a Trackello type using your preferences application token and appkey.
func NewTrackello(token, appKey string) (*Trackello, error) {

	if len(token) == 0 || len(appKey) == 0 {
		return nil, fmt.Errorf("Token and/or AppKey is nil: Token: %s, AppKey: %s", token, appKey)
	}
	// New Trello Client
	tr, err := trello.NewAuthClient(appKey, &token)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Trackello{
		token:  token,
		appKey: appKey,
		client: tr,
	}, nil
}

// Board will return the Trello Board given it's ID string.
func (t *Trackello) Board(id string) (*trello.Board, error) {
	board, err := t.client.Board(id)
	if err != nil {
		return board, fmt.Errorf("Error retrieving Board ID '%s': %v", id, err)
	}
	return board, nil
}

// Boards will list all of the boards for the authenticated user (i.e. 'me').
func (t *Trackello) Boards() ([]trello.Board, error) {
	member, err := t.client.Member("me")
	if err != nil {
		log.Fatalf("Error getting 'me' Member: %v", err)
		return make([]trello.Board, 0), err
	}
	boards, err := member.Boards()
	return boards, err
}
