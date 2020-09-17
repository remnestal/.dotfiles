package aggregated

import (
	structs "gitlab.com/abios/v3-structs"
	"gitlab.com/abios/v3-structs/pbp/cs"
)

type DamageBreakdown struct {
	cs.Hitpoints
	Hitgroups struct {
		Head uint16 `json:"head"`
	} `json:"hitgroups"`
}

type PlayerStats struct {
	Kills struct {
		Total     uint16 `json:"total"`
		Headshots uint16 `json:"headshots"`
		WallBangs uint16 `json:"wallbangs"`
	} `json:"kills"`

	Assists struct {
		Total uint16 `json:"total"`
		Flash uint16 `json:"flash"`
	} `json:"assists"`

	Bombs struct {
		Planted uint16 `json:"planted"`
		Defused uint16 `json:"defused"`
	} `json:"bombs"`

	Damage struct {
		Given struct {
			DamageBreakdown
			PerWeapon []struct {
				Weapon structs.GameAsset `json:"weapon"`
				DamageBreakdown
			} `json:"per_weapon"`
		} `json:"dealt"`
		Taken DamageBreakdown `json:"taken"`
	} `json:"damage"`
}

type Player struct {
	structs.Edge

	Match struct {
		PlayerStats
		Deaths uint16 `json:"deaths"`
	} `json:"match_stats"`

	Round struct {
		PlayerStats
		HitPoints cs.Hitpoints `json:"hitpoints"`
	} `json:"round_stats"`

	NetWorth struct {
		Money          uint16 `json:"money"`
		EquipmentValue uint16 `json:"equipment_value"`
	} `json:"net_worth"`

	Inventory cs.Inventory `json:"inventory"`
}
