package events

type KilledEvent struct {
	AttackerAbstraction
	VictimAbstraction
	Headshot    bool `json:"headshot"`
	IndirectHit bool `json:"indirect_hit"`
}

type KilledSelfEvent struct {
	AttackerAbstraction
}

type KilledByBombEvent struct {
	BombEventMetadata
}
