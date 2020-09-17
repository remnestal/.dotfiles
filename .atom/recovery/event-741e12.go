package cs

import (
	"time"

	"gitlab.com/abios/motattack/pkg/csgods/events"
	structs "gitlab.com/abios/v3-structs"
)

/*
	TODO: We can only really re-export motattack eventypes if both (or neither)
	are open source
*/

type EventType string

const (
	StateEventType EventType = "state"
	PhaseEventType EventType = "phase"

	GameOverEventType            EventType = events.GameOverEventType
	MolotovEventType             EventType = events.MolotovEventType
	PlayerAttackedEventType      EventType = events.PlayerAttackedEventType
	PlayerKilledPlayerEventType  EventType = events.PlayerKilledPlayerEventType
	PlayerDestroyedEventType     EventType = events.PlayerDestroyedEventType
	PlayerSwitchedEventType      EventType = events.PlayerSwitchedEventType
	PlayerPickedUpEventType      EventType = events.PlayerPickedUpEventType
	PlayerPurchasedEventType     EventType = events.PlayerPurchasedEventType
	PlayerDroppedEventType       EventType = events.PlayerDroppedEventType
	PlayerLeftBuyzoneEventType   EventType = events.PlayerLeftBuyzoneEventType
	PlayerHasEventType           EventType = events.PlayerHasEventType
	PlayerChatEventType          EventType = events.PlayerChatEventType
	PlayerAquiredBombEventType   EventType = events.PlayerAquiredBombEventType
	PlayerDroppedBombEventType   EventType = events.PlayerDroppedBombEventType
	PlayerPlantedBombEventType   EventType = events.PlayerPlantedBombEventType
	PlayerDefusingBombEventType  EventType = events.PlayerDefusingBombEventType
	PlayerDefusedBombEventType   EventType = events.PlayerDefusedBombEventType
	PlayerKilledByBombEventType  EventType = events.PlayerKilledByBombEventType
	PlayerSuicideEventType       EventType = events.PlayerSuicideEventType
	PlayerAssistedEventType      EventType = events.PlayerAssistedEventType
	PlayerBlindedEventType       EventType = events.PlayerBlindedEventType
	PlayerEconomyChangeEventType EventType = events.PlayerEconomyChangeEventType
	PlayerActivatedEventType     EventType = "activated" // activated grenade
	PlayerThrewEventType         EventType = events.PlayerThrewEventType
	StartedMapEventType          EventType = events.StartedMapEventType
	StartFreezetimeEventType     EventType = events.StartFreezetimeEventType
	TeamWonEventType             EventType = events.TeamWonEventType
	TeamFactionEventType         EventType = "team-faction"
	MatchStartedEventType        EventType = events.MatchPausedEventType
	MatchPausedEventType         EventType = events.MatchPausedEventType
	MatchResumedEventType        EventType = events.MatchResumedEventType
	RoundStartedEventType        EventType = events.RoundStartedEventType
	RoundRestartEventType        EventType = events.RoundRestartEventType
	VoteStartedEventType         EventType = events.VoteStartedEventType
	VoteCastEventType            EventType = events.VoteCastEventType
	VoteFailedEventType          EventType = events.VoteFailedEventType
	VoteSucceededEventType       EventType = events.VoteSucceededEventType
	AccoladeEventType            EventType = events.AccoladeEventType
)

type Event struct {
	Type         EventType     `json:"type"`
	Timestamp    time.Time     `json:"timestamp"`
	Match        structs.Match `json:"match"`
	CurrentRound uint          `json:"current_round"`
	Payload      interface{}   `json:"payload"`
}
