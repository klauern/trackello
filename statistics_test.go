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
