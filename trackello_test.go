package trackello

import (
	"os"
	"testing"
)

func setupTrackello(t *testing.T) *Trackello {
	client, err := NewTrackello(os.Getenv(TRACKELLO_TOKEN), os.Getenv(TRACKELLO_APPKEY))
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestBoard(t *testing.T) {
	var boardTests = []struct {
		boardID  string
		hasError bool
	}{
		{"test", true},
		{"5532c8c02c1b8cbebb3e4de5", false},
		{"550ce6ae4285507e2c51f661", false},
		{"550ce6ae4285507e2c", true},
	}
	client := setupTrackello(t)
	for _, tt := range boardTests {
		_, e := client.Board(tt.boardID)
		hasErr := (e != nil)
		if tt.hasError != hasErr {
			t.Fatalf("Expected %t for boardID '%s', got %t", tt.hasError, tt.boardID, hasErr)
		}
	}
}

func TestListBoards(t *testing.T) {
	client := setupTrackello(t)
	boards, err := client.Boards()
	if err != nil {
		t.Fatalf("Error getting boards: %v", err)
	}
	if len(boards) < 2 {
		t.Fatalf("Not enough boards.  Expected %v, got %v boards", 2, len(boards))
	}
}

func TestListBoardsWithBadIds(t *testing.T) {
	var tokenKeyTests = []struct {
		token  string
		appkey string
	}{
		//	{"", ""},
		//{"", os.Getenv(TRACKELLO_APPKEY)},
		//{os.Getenv(TRACKELLO_TOKEN), ""},
		{os.Getenv(TRACKELLO_TOKEN), "JUNKID"},
		{"JUNKID", os.Getenv(TRACKELLO_APPKEY)},
	}
	for _, tt := range tokenKeyTests {
		track, _ := NewTrackello(tt.token, tt.appkey)
		_, err := track.Board("")
		if err == nil {
			t.Fatalf("Expected an error with a bad Board ID, did NOT get one")
		}
		_, err = track.Board("JUNKID")
		if err == nil {
			t.Fatalf("Expected an error with a bad Board ID, did NOT get one")
		}
	}

}

func TestNoTokenOrAppKey(t *testing.T) {
	var tokenKeyTests = []struct {
		token  string
		appkey string
	}{
		{"", ""},
		{"", os.Getenv(TRACKELLO_APPKEY)},
		{os.Getenv(TRACKELLO_TOKEN), ""},
		//{os.Getenv(TRACKELLO_TOKEN), "JUNKID"},
		//{"JUNKID", os.Getenv(TRACKELLO_APPKEY)},
	}
	for _, tt := range tokenKeyTests {
		_, err := NewTrackello(tt.token, tt.appkey)
		if err == nil {
			t.Fatalf("Expected an error with %s for a token and %s for an AppKey.  Did NOT get one", tt.token, tt.appkey)
		}
	}
}

func TestBadBoardId(t *testing.T) {
	client := setupTrackello(t)
	_, err := client.Board("")
	if err == nil {
		t.Fatalf("Expected an error for an empty BoardID.  Did NOT get an error")
	}
	_, err = client.Board("JUNKID")
	if err == nil {
		t.Fatalf("Expected an error for a junk BoardID.  Did NOT get an error")
	}
}
