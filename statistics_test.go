package trackello

import "testing"

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
