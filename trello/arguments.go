package trello

import (
	"time"

	"github.com/VojtechVitek/go-trello"
)

const (
	DateLayout      string = "2006-01-02T15:04:05Z"
	API_APPKEY      string = "TRELLO_APPKEY"
	API_TOKEN       string = "TRELLO_TOKEN"
	PREFERRED_BOARD string = "preferredBoard"
)

func CreateArgsForBoardActions() []*trello.Argument {
	var args []*trello.Argument
	twoWeeksAgo := time.Now().Add(-1 * time.Hour * 24 * 14).Format(DateLayout)
	args = append(args, trello.NewArgument("since", twoWeeksAgo))
	args = append(args, trello.NewArgument("limit", "500"))
	return args
}
