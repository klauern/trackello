package trackello

import (
	"os"
	"testing"
)

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
	client, err := NewTrackello(os.Getenv(TRACKELLO_TOKEN), os.Getenv(TRACKELLO_APPKEY))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range boardTests {
		_, e := client.Board(tt.boardID)
		hasErr := (e != nil)
		if tt.hasError != hasErr {
			t.Fatalf("Expected %t for boardID '%s', got %t", tt.hasError, tt.boardID, hasErr)
		}
	}
}
