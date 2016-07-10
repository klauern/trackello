package trackello

// cardStatistics provides a way to show various pieces of
type cardStatistics struct {
	comments, // represented by a horizontal ellepsis ⋯ 0x22EF
	updates, // represented by a keyboard 0x2328
	checklistsCreated, // represented by plus +
	checklistItemsChecked int // represented by check mark ✓ 0x2713
}

type listStatistics cardStatistics

type Statistics interface {
	Calculate()
	Print() string
}

func (s *listStatistics) Calculate() {}

func (s *listStatistics) Print() string {
	return ""
}

func (s *cardStatistics) Calculate() {}

func (s *cardStatistics) Print() string {
	return ""
}
