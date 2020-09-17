package passthrough

import structs "gitlab.com/abios/v3-structs"

type KillEvent struct {
	ConfrontationMetadata
	Weapon      structs.GameAsset `json:"weapon"`
	Headshot    bool              `json:"headshot"`
	IndirectHit bool              `json:"indirect_hit"`
}
