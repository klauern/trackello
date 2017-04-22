package trackello

import (
	"os"
	"testing"
)

func SetupClient(t *testing.T) *Trackello {
	c, err := NewTrackello(os.Getenv(TRACKELLO_TOKEN), os.Getenv(TRACKELLO_APPKEY))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func SetupTest(t *testing.T, id string) *Board {
	c := SetupClient(t)
	b, err := c.Board(id)
	if err != nil {
		t.Fatal(err)
	}
	return NewBoard(b)
}

func TestPopulateList(t *testing.T) {
	boardIDs := []string{"550ce6ae4285507e2c51f661", "56f5588daa643d38a1c7f111", "57269197fc77edb4599f2383"}
	for _, id := range boardIDs {
		board := SetupTest(t, id)
		board.PopulateLists()
	}

}

func TestMapActions(t *testing.T) {
	boardIDs := []string{"550ce6ae4285507e2c51f661", "56f5588daa643d38a1c7f111", "57269197fc77edb4599f2383"}
	for _, id := range boardIDs {
		board := SetupTest(t, id)
		board.MapActions()

	}
}

func TestPrintActions(t *testing.T) {
	boardIDs := []string{"550ce6ae4285507e2c51f661", "56f5588daa643d38a1c7f111", "57269197fc77edb4599f2383"}
	for _, id := range boardIDs {
		board := SetupTest(t, id)
		board.PrintActions()
	}
}

func TestBadBoard(t *testing.T){
	client := SetupClient(t)
	b, err := client.Board("BADID")
	if err == nil{
		t.Errorf("Expected Error, Got %v", b)
	}
}