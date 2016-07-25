package trackello

import (
	"testing"
	"os"
)

func SetupTest(t *testing.T) *Board {
	c, err := NewTrackello(os.Getenv(TRACKELLO_TOKEN), os.Getenv(TRACKELLO_APPKEY))
	if err != nil {
		t.Fatal(err)
	}
	b, err := c.Board("550ce6ae4285507e2c51f661")
	if err != nil {
		t.Fatal(err)
	}
	return NewBoard(b)
}

func TestPopulateList(t *testing.T) {
	board := SetupTest(t)
	board.PopulateLists()

}

func TestMapActions(t *testing.T) {
	board := SetupTest(t)
	board.MapActions()
}

func TestPrintActions(t *testing.T) {
	board := SetupTest(t)
	                  board.PrintActions()
}
