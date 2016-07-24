package trackello

import (
	"fmt"
	"os"
	"testing"
)

func TestBoard(t *testing.T) {
	var boardTests = []struct {
		boardId  string
		hasError bool
	}{
		{"test", true},
		{"5532c8c02c1b8cbebb3e4de5", false},
		{"550ce6ae4285507e2c51f661", false},
		{"550ce6ae4285507e2c", true},
	}
	client, err := NewTrackello(os.Getenv("TRACKLLEO_TOKEN"), os.Getenv("TRACKELLO_APPKEY"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range boardTests {
		_, e := client.Board(tt.boardId)
		hasErr := (e != nil)
		if hasErr {
			fmt.Printf("Error %v\n", e)
		}
		if tt.hasError != hasErr {
			t.Fatalf("Expected %t for boardID '%s', got %t", tt.hasError, tt.boardId, hasErr)
		}
	}
}
