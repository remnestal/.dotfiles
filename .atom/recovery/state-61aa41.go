package state

import (
	structs "gitlab.com/abios/v3-structs"
	"gitlab.com/abios/v3-structs/pbp/cs"
)

type State struct {
	Map         structs.Map `json:"map"`
	MatchLength uint64      `json:"match_length"`
	RoundNumber uint64      `json:"round_number"`

	MatchPhase cs.Phase `json:"match_phase"`
	RoundPhase cs.Phase `json:"round_phase"`

	Teams struct {
		Home Team `json:"home"`
		Away Team `json:"away"`
	} `json:"teams"`

	PastRounds []Round `json:"past_rounds"`
}
