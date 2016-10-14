package trackello

import (
	"testing"

	trello "github.com/VojtechVitek/go-trello"
)

func TestTotalStatisticsCount(t *testing.T) {
	stats := &Statistics{
		comments:             10,
		checkListItemUpdates: 10,
		creates:              10,
		updates:              10,
	}

	if stats.Total() != 40 {
		t.Fatalf("Expected %d, got %d", 40, stats.Total())
	}
}

func TestPrintStatistics(t *testing.T) {
	stats := &Statistics{
		comments:             10,
		checkListItemUpdates: 10,
		creates:              10,
		updates:              10,
	}

	printed := stats.PrintStatistics()

	//expected := "[\x1b[96m10 +\x1b[0m\x1b[91m 10 ≡\x1b[0m\x1b[92m 10 ✓\x1b[0m\x1b[95m 10 …\x1b[0m]"
	expected := "[10 + 10 ≡ 10 ✓ 10 …]"
	if printed != expected {
		t.Fatalf("Printed output doesn't match expected: got %s, wanted %s", printed, expected)
	}
}

func TestPrintEmptyStatistics(t *testing.T) {
	stats := &Statistics{}
	printed := stats.PrintStatistics()
	expected := "[0  + 0  ≡ 0  ✓ 0  …]"
	if printed != expected {
		t.Fatalf("TestPrintEmptyStatistics Failed: expected %s, got %s", expected, printed)
	}
}

func TestAddStatistics(t *testing.T) {
	tests := []struct {
		statType string
		comments int
		updates  int
		creates  int
		checks   int
	}{
		{"commentCard", 1, 0, 0, 0},
		{"updateCard", 0, 1, 0, 0},
		{"deleteAttachmentFromCard", 0, 1, 0, 0},
		{"updateList", 0, 1, 0, 0},
		{"createCard", 0, 0, 1, 0},
		{"addChecklistToCard", 0, 0, 1, 0},
		{"addAttachmentToCard", 0, 0, 1, 0},
		{"updateCheckItemStateOnCard", 0, 0, 0, 1},
		//{"junkItemAction", 0, 0, 0, 1},
	}
	for _, tt := range tests {
		card := NewCard(trello.Card{})
		err := card.AddCalculation(trello.Action{
			Type: tt.statType,
		})
		if err != nil {
			t.Fatalf("Error when adding a calculation of type %s: %v", tt.statType, err)
		}
		stats := card.stats
		if stats.checkListItemUpdates != tt.checks {
			t.Fatalf("Error with 'checkListItemUpdates' for statType %s.  Expected %v, got %v", tt.statType, tt.checks, stats.checkListItemUpdates)
		}
		if stats.comments != tt.comments {
			t.Fatalf("Error with 'comments' for statType %s.  Expected %v, got %v", tt.statType, tt.comments, stats.comments)
		}
		if stats.creates != tt.creates {
			t.Fatalf("Error with 'creates' for statType %s.  Expected %v, got %v", tt.statType, tt.creates, stats.creates)
		}
		if stats.updates != tt.updates {
			t.Fatalf("Error with 'updates' for statType %s.  Expected %v, got %v", tt.statType, tt.updates, stats.updates)
		}
	}
}

func TestErrorStatisticsCases(t *testing.T) {
	card := NewCard(trello.Card{})
	err := card.AddCalculation(trello.Action{Type: "JunkStatusType"})
	if err == nil {
		t.Fatalf("Expected an error for an unknown type. Got no error")
	}
	if card.stats.updates != 1 {
		t.Fatalf("expected card statistics to have %v for updates.  Got %v instead", 1, card.stats.updates)
	}
}
