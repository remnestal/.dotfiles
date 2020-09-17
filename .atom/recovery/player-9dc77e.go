package events

const (
	PlayerAttackedEventType          EventType = "attacked"
	PlayerKilledPlayerEventType      EventType = "killed"
	PlayerDestroyedEventType         EventType = "destroyed"
	PlayerValidatedIdEventType       EventType = "validated-id"
	PlayerEnteredGameEventType       EventType = "entered-game"
	PlayerEnteredNoclipModeEventType EventType = "entered-noclip-mode"
	PlayerConnectedEventType         EventType = "connected"
	PlayerDisconnectedEventType      EventType = "disconnected"
	PlayerChangedNameEventType       EventType = "changed-name"
	PlayerSwitchedEventType          EventType = "switched-faction"
	PlayerPickedUpEventType          EventType = "picked-up"
	PlayerPurchasedEventType         EventType = "purchased"
	PlayerDroppedEventType           EventType = "dropped"
	PlayerLeftBuyzoneEventType       EventType = "left-buyzone"
	PlayerAcquiredBombEventType      EventType = "aquired-bomb"
	PlayerDroppedBombEventType       EventType = "dropped-bomb"
	PlayerPlantedBombEventType       EventType = "planted-bomb"
	PlayerDefusingBombEventType      EventType = "defusing-bomb"
	PlayerDefusedBombEventType       EventType = "defused-bomb"
	PlayerKilledByBombEventType      EventType = "killed-by-bomb"
	PlayerSuicideEventType           EventType = "suicide"
	PlayerAssistedEventType          EventType = "assisted"
	PlayerBlindedEventType           EventType = "blinded"
	PlayerEconomyChangedEventType    EventType = "economy-change"
	PlayerHasEventType               EventType = "has"
	PlayerThrewEventType             EventType = "threw"
	PlayerChatEventType              EventType = "chat"
)

type Player struct {
	Nickname string  `json:"nickname"`
	LocalId  float64 `json:"local_id"`
	SteamId  string  `json:"steam_id"`
	Side     string  `json:"side"`
	Position *Vector `json:"position"`
}

func (src Player) Copy() (dst Player) {
	dst = src
	if src.Position != nil {
		pos := *src.Position
		dst.Position = &pos
	}
	return
}

type PlayerEventData struct {
	E
	Player Player `json:"player"`
}

func (src PlayerEventData) Copy() (dst PlayerEventData) {
	dst = src
	dst.E = *src.E.Copy().(*E)
	dst.Player = src.Player.Copy()
	return
}

type DamageBreakdown struct {
	Health float64 `json:"health"`
	Armor  float64 `json:"armor"`
}

type PlayerAttackedEvent struct {
	PlayerEventData
	With      Equipment       `json:"with"`
	Victim    Player          `json:"victim"`
	Where     Hitgroup        `json:"where"`
	Damage    DamageBreakdown `json:"damage"`
	Remaining DamageBreakdown `json:"remaining"`
}

func (src PlayerAttackedEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	dst.Victim = src.Victim.Copy()
	return &dst
}

type PlayerKilledEventData struct {
	PlayerEventData
	With          Equipment `json:"with"`
	Penetrated    bool      `json:"penetrated"`
	AttackerBlind bool      `json:"attacker_blind"`
	Scoped        *bool     `json:"scoped"`
}

func (src PlayerKilledEventData) Copy() (dst PlayerKilledEventData) {
	dst = src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	if src.Scoped != nil {
		pos := *src.Scoped
		dst.Scoped = &pos
	}
	return
}

type PlayerKilledPlayerEvent struct {
	PlayerKilledEventData
	Victim       Player `json:"victim"`
	Headshot     bool   `json:"headshot"`
	ThroughSmoke bool   `json:"through_smoke"`
}

func (src PlayerKilledPlayerEvent) Copy() Event {
	dst := src
	dst.PlayerKilledEventData = src.PlayerKilledEventData.Copy()
	dst.Victim = src.Victim.Copy()
	return &dst
}

type PlayerDestroyedEvent struct {
	PlayerKilledEventData
	What  string `json:"what"`
	Where Vector `json:"at"`
}

func (src PlayerDestroyedEvent) Copy() Event {
	dst := src
	dst.PlayerKilledEventData = src.PlayerKilledEventData.Copy()
	return &dst
}

type PlayerValidatedIdEvent struct {
	PlayerEventData
}

func (src PlayerValidatedIdEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerEnteredGameEvent struct {
	PlayerEventData
}

type PlayerEnteredNoclipModeEvent struct {
	PlayerEventData
}

func (src PlayerEnteredGameEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerConnectedEvent struct {
	PlayerEventData
	From string `json:"from"`
}

func (src PlayerConnectedEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerDisconnectedEvent struct {
	PlayerEventData
	Reason string `json:"reason"`
}

func (src PlayerDisconnectedEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerChangedNameEvent struct {
	PlayerEventData
	Previously string `json:"previously"`
}

func (src PlayerChangedNameEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerSwitchedEvent struct {
	PlayerEventData
	From string `json:"from"`
	To   string `json:"to"`
}

func (src PlayerSwitchedEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerPickedUpEvent struct {
	PlayerEventData
	What Equipment `json:"what"`
}

func (src PlayerPickedUpEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerPurchasedEvent struct {
	PlayerEventData
	What Equipment `json:"what"`
}

func (src PlayerPurchasedEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerDroppedEvent struct {
	PlayerEventData
	What Equipment `json:"what"`
}

func (src PlayerDroppedEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerLeftBuyzoneEvent struct {
	PlayerEventData
	With []Equipment `json:"with"`
}

func (src PlayerLeftBuyzoneEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	copy(dst.With, src.With)
	return &dst
}

type PlayerAcquiredBombEvent struct {
	PlayerEventData
}

func (src PlayerAcquiredBombEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerDroppedBombEvent struct {
	PlayerEventData
}

func (src PlayerDroppedBombEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerPlantedBombEvent struct {
	PlayerEventData
}

func (src PlayerPlantedBombEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerDefusingBombEvent struct {
	PlayerEventData
	Kit bool `json:"kit"`
}

func (src PlayerDefusingBombEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerDefusedBombEvent struct {
	PlayerEventData
}

func (src PlayerDefusedBombEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerKilledByBombEvent struct {
	PlayerEventData
}

func (src PlayerKilledByBombEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerSuicideEvent struct {
	PlayerEventData
	With Equipment `json:"with"`
}

func (src PlayerSuicideEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerAssistedEvent struct {
	PlayerEventData
	Flash  bool   `json:"flash"`
	Victim Player `json:"victim"`
}

func (src PlayerAssistedEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	dst.Victim = src.Victim.Copy()
	return &dst
}

type PlayerBlindedEvent struct {
	PlayerEventData
	By       Player  `json:"by"`
	For      float64 `json:"for"`
	Entindex float64 `json:"entindex"`
}

func (src PlayerBlindedEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	dst.By = src.By.Copy()
	return &dst
}

type PlayerEconomyChangedEvent struct {
	PlayerEventData
	Before    float64   `json:"before"`
	After     float64   `json:"after"`
	Delta     float64   `json:"delta"`
	Equipment Equipment `json:"equipment,omitempty"`
}

func (src PlayerEconomyChangedEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerHasEvent struct {
	PlayerEventData
	Money float64 `json:"money"`
}

func (src PlayerHasEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerThrewEvent struct {
	PlayerEventData
	What     Equipment `json:"what"`
	Position Vector    `json:"position"`
	Entindex float64   `json:"entindex"`
}

func (src PlayerThrewEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}

type PlayerChatEvent struct {
	PlayerEventData
	Receivers string `json:"receivers"`
	Message   string `json:"message"`
}

func (src PlayerChatEvent) Copy() Event {
	dst := src
	dst.PlayerEventData = src.PlayerEventData.Copy()
	return &dst
}
