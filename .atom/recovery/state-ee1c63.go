package state

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"gitlab.com/abios/motattack/pkg/csgods"
	"gitlab.com/abios/motattack/pkg/csgods/events"
	"gitlab.com/abios/motattack/pkg/csgods/events/state/strutil"
	"gitlab.com/abios/v3-structs/pbp/cs"
)

const (
	FULL_HEALTH = 100
	FULL_ARMOR  = 100

	DEFAULT_PISTOL_COST float64 = 500
	DEFAULT_KEVLAR_COST float64 = 650
	DEFAULT_HELMET_COST float64 = 350

	EQUIPMENT_ID_KNIFE       = "knife"
	EQUIPMENT_ID_BAYONET     = "bayonet"
	EQUIPMENT_ID_VESTHELM    = "vesthelm"
	EQUIPMENT_ID_VEST        = "vest"
	EQUIPMENT_ID_HELMET      = "helmet"
	EQUIPMENT_ID_KEVLAR      = "kevlar"
	EQUIPMENT_ID_FLASHBANG   = "flashbang"
	EQUIPMENT_ID_UNWANTED_C4 = "C4"

	PLAYER_TYPE_BOT = "BOT"

	CVAR_STARTING_MONEY = "mp_startmoney"
	CVAR_MAX_MONEY      = "mp_maxmoney"
)

type State struct {
	Map string `json:"map"`

	/*
		Phase represents the current phase of the match and can be
		*) warmup
		*) live
		*) intermission
		*) gameover
	*/
	Phase  MatchPhase `json:"phase"`
	Length uint       `json:"length_minutes"`

	// Round represents the state of the current round being played
	Round struct {
		Number uint `json:"number"`

		/*
			Phase describes the current phase of the round.
			Valid values are:
			*) live
			*) over
			*) freezetime
			*) warmup
			*) bomb
			*) defuse
		*/
		Phase RoundPhase `json:"phase"`
	} `json:"round"`

	// Match session variables
	Cvars Cvars `json:"cvars"`

	Home Team `json:"home"`
	Away Team `json:"away"`

	// Rounds holds the conlusion of all previosly played rounds
	Rounds []Round `json:"rounds"`

	equipmentCost       map[string]float64
	playerSwitchCounter uint8
}

type Cvars struct {
	StartMoney float64 `json:"mp_startmoney"`
	MaxMoney   float64 `json:"mp_maxmoney"`
}

type Round struct {
	Number uint   `json:"number"`
	Winner string `json:"winner"`
	Reason string `json:"reason"`
}

type Team struct {
	Score   float64  `json:"score"`
	Faction string   `json:"side"`
	Players []Player `json:"players"`
}

type Player struct {
	Nickname string `json:"nickname"`
	SteamId  string `json:"steam_id"`

	// Stats is the overall stats of the match
	Stats struct {
		Kills   float64 `json:"kills"`
		Assists float64 `json:"assists"`
		Deaths  float64 `json:"deaths"`
		Adr     float64 `json:"adr"`
	} `json:"stats"`

	// State is the player's statistics in this round
	State struct {
		Health float64 `json:"health"`

		// Current armor 0-100
		Armor  float64 `json:"armor"`
		Helmet bool    `json:"helmet"`

		Money      float64 `json:"money"`
		EquipValue float64 `json:"equipment_value"`

		Kills struct {
			Bodyshot float64 `json:"bodyshot"`
			Headshot float64 `json:"headshot"`
		} `json:"kills"`

		Damage struct {
			Health float64 `json:"health"`
			Armor  float64 `json:"armor"`
		} `json:"damage"`
	} `json:"state"`

	Inventory cs.Inventory `json:"inventory"`
	Equipment []string     `json:"equipment"`
}

func (p1 Player) equals(p2 events.Player) bool {
	if p1.SteamId == PLAYER_TYPE_BOT && p2.SteamId == PLAYER_TYPE_BOT {
		return p1.Nickname == p2.Nickname
	} else {
		return p1.SteamId == p2.SteamId
	}
}

func New(cvars *Cvars) *State {
	s := State{
		Home: Team{
			Faction: "ct",
		},
		Away: Team{
			Faction: "t",
		},
	}

	s.Map = "undecided"
	s.equipmentCost = map[string]float64{
		"glock":        DEFAULT_PISTOL_COST,
		"hkp2000":      DEFAULT_PISTOL_COST,
		"usp_silencer": DEFAULT_PISTOL_COST,
	}

	// Reset player states and game phases
	s.Reset()
	s.Phase.Reset()
	s.Round.Phase.Reset()

	// Set any predefined cvars
	if cvars != nil {
		s.Cvars = *cvars
	}
	return &s
}

func (s *State) resetPlayer(p Player) Player {
	p.Stats.Kills = 0
	p.Stats.Assists = 0
	p.Stats.Deaths = 0
	p.Stats.Adr = 0

	p.State.Health = FULL_HEALTH
	p.State.Armor = 0
	p.State.Helmet = false
	p.State.Money = s.Cvars.StartMoney

	p.State.Kills.Bodyshot = 0
	p.State.Kills.Headshot = 0

	p.State.Damage.Health = 0
	p.State.Damage.Armor = 0

	p.State.EquipValue = 0
	p.Equipment = []string(nil)
	return p
}

func (s *State) switchSides() {
	s.Home.Faction, s.Away.Faction = s.Away.Faction, s.Home.Faction
	s.Home.Score, s.Away.Score = s.Away.Score, s.Home.Score
}

func (s *State) Reset() {
	// Reset the list of past rounds
	s.Round.Number = 0
	s.Rounds = []Round{}
	// Reset both teams
	for _, t := range []*Team{&s.Home, &s.Away} {
		t.Score = 0
		// Reset all stats, state and equipment for alla players
		for i := 0; i < len(t.Players); i++ {
			t.Players[i] = s.resetPlayer(t.Players[i])
		}
	}
}

var ErrCvarNotRecognized = errors.New("Not a recognized Cvar key-value pair")

func (s *State) setCvar(key, value string) error {
	switch key {
	case CVAR_STARTING_MONEY:
		if f, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("Unable to parse startmoney: %w", err)
		} else {
			s.Cvars.StartMoney = f
		}

	case CVAR_MAX_MONEY:
		if f, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("Unable to parse maxmoney: %w", err)
		} else {
			s.Cvars.MaxMoney = f
		}

	default:
		return ErrCvarNotRecognized
	}
	return nil
}

func (t *Team) EveryPlayer(condition func(Player) bool) bool {
	for _, p := range t.Players {
		if !condition(p) {
			return false
		}
	}
	return true
}

func (s *State) updateSide(side string, transform func(Team) Team) {
	switch side {
	case s.Home.Faction:
		s.Home = transform(s.Home)
	case s.Away.Faction:
		s.Away = transform(s.Away)
	}
}

func (s *State) updateSides(transform func(Team) Team) {
	s.Home = transform(s.Home)
	s.Away = transform(s.Away)
}

func (s *State) updatePlayer(player events.Player, transform func(Player) Player) {
	s.updateSide(player.Side, func(current Team) Team {
		found := false
		for i, p := range current.Players {
			if p.equals(player) {
				current.Players[i] = transform(p)
				found = true
				break
			}
		}

		if !found {
			current.Players = append(current.Players, transform(s.resetPlayer(Player{
				Nickname: player.Nickname,
				SteamId:  player.SteamId,
			})))
		}
		return current
	})
}

func (s *State) Update(_event events.Event) {
	switch event := _event.(type) {

	case *events.Cvar:
		if err := s.setCvar(event.Cvar, event.Value); err != ErrCvarNotRecognized && err != nil {
			log.Fatalf("Unable to set Cvar %v: %v", event.Cvar, err)
		}

	case *events.RconCommands:
		for _, cmd := range event.Commands {
			components := strings.Split(cmd, " ")
			if len(components) == 2 {
				if err := s.setCvar(components[0], components[1]); err != ErrCvarNotRecognized && err != nil {
					log.Fatalf("Unable to set Cvar %v: %v", components[0], err)
				}
			}
		}

	case *events.PlayerDisconnected:
		s.updateSide(event.Player.Side, func(current Team) Team {
			for i, p := range current.Players {
				if p.equals(event.Player) {
					current.Players = append(current.Players[:i], current.Players[i+1:]...)
					break
				}
			}
			return current
		})

	case *events.PlayerSwitched:
		if s.Phase.In(PHASE_LIVE) {
			s.playerSwitchCounter++
			if s.playerSwitchCounter == 10 {
				s.playerSwitchCounter = 0
				s.switchSides()
			}
		}
		// Remove from both sides.
		//
		// Conceptually a player should only be part of one side, but somehow,
		// sometimes, players switch from "unassigned" when they are in fact already
		// part of the side they're joining. Judging by the test data on which this
		// code is developed with regards to, the player should be reset in either
		// case. Thus, removing from both teams is a safety precaution.
		s.updateSides(func(current Team) Team {
			for i, p := range current.Players {
				if p.equals(event.Player) {
					current.Players = append(current.Players[:i], current.Players[i+1:]...)
					break
				}
			}
			return current
		})
		// Add to new side
		s.updateSide(event.To, func(current Team) Team {
			current.Players = append(current.Players, s.resetPlayer(Player{
				Nickname: event.Player.Nickname,
				SteamId:  event.Player.SteamId,
			}))
			return current
		})

	case *events.PlayerPickedUp:
		s.updatePlayer(event.Player, func(current Player) Player {
			if event.What == EQUIPMENT_ID_VESTHELM {
				// Special case if vesthelm is picked up; this is a composite
				// item consisting of kevlar(100) and helmet
				current.Equipment = strutil.Remove(current.Equipment, func(item string) bool {
					return item == EQUIPMENT_ID_HELMET || item == EQUIPMENT_ID_KEVLAR
				})
				current.Equipment = append(current.Equipment, EQUIPMENT_ID_HELMET, EQUIPMENT_ID_KEVLAR)
				current.State.Helmet = true
				current.State.Armor = FULL_ARMOR

			} else if event.What == EQUIPMENT_ID_VEST {
				// Add the armor value represented in the kevlar entity ID
				// but simply add "kevlar" to the list of equipment
				current.Equipment = strutil.Remove(current.Equipment, func(item string) bool {
					return item == EQUIPMENT_ID_KEVLAR
				})
				current.Equipment = append(current.Equipment, event.What)
				current.State.Armor = FULL_ARMOR

			} else {
				// Otherwise add the item as-is
				current.Equipment = append(current.Equipment, event.What)
			}
			return current
		})

	case *events.PlayerDropped:
		s.updatePlayer(event.Player, func(current Player) Player {

			// Normalize knife IDs by attempting to parse the item as a
			// knife ID
			if ok, _ := parseKnife(event.What); ok {
				current.Equipment = strutil.Remove(current.Equipment, func(item string) bool {
					ok, _ := parseKnife(item)
					return ok
				})

			} else {
				// Otherwise only strutil.Remove the specified item (once)
				current.Equipment = strutil.Remove(current.Equipment, func(item string) bool {
					return item == event.What
				})
			}
			return current
		})

	case *events.StartFreezetime:
		switch {
		case s.Phase.In(PHASE_PREGAME):
			// This is required to adequately "clear" all player's inventories before
			// the knife-round in some games. This fix might, however, be a
			// coincidence and clearing the state on every freezetime during pregame
			// makes the state pretty worthless until the game goes live/knife
			s.Reset()

		case s.Phase.In(PHASE_LIVE):
			for _, team := range []*Team{&s.Home, &s.Away} {
				s.updateSide(team.Faction, func(t Team) Team {
					for i := range t.Players {
						t.Players[i].State.Health = FULL_HEALTH
					}
					return t
				})
			}
		}

	case *events.StartedMap:
		*s = *New(&s.Cvars)
		s.Map = event.Map

	case *events.TeamWon:
		s.updateSide(csgods.TERRORISTS, func(t Team) Team {
			t.Score = event.Score.T
			return t
		})
		s.updateSide(csgods.COUNTER_TERRORISTS, func(t Team) Team {
			t.Score = event.Score.CT
			return t
		})
		// Update list of past rounds
		s.Rounds = append(s.Rounds, Round{
			Number: s.Round.Number,
			Winner: event.TeamEvent.Team,
			Reason: event.Reason,
		})

	case *events.PlayerEconomyChange:
		s.updatePlayer(event.Player, func(p Player) Player {
			// Register the cost of the item if it is not already noted
			if _, exist := s.equipmentCost[event.Equipment]; !exist {
				s.equipmentCost[event.Equipment] = event.Delta
			}
			// Sanity check of the player's economy prior to the change
			if p.State.Money != event.Before {
				// FIXME: how should this be handled/logged?
				log.Printf("The economy tracking of player `%v` is wrong: %v != %v\n", p.Nickname, p.State.Money, event.Before)
			}
			// Alter the player's economy
			p.State.Money = event.After
			return p
		})

	case *events.PlayerLeftBuyzone:
		s.updatePlayer(event.Player, func(p Player) Player {

			// The bomb is represented as both `weapon_c4` and plain `C4`
			// in the buyzone equipment summary. Remove the latter.
			inventory := strutil.Remove(event.With, func(item string) bool {
				return item == EQUIPMENT_ID_UNWANTED_C4
			})

			// Set the armor-value represented by the kevlar item, but
			// simply list it as `kevlar` and not `kevlar(x)`
			inventory = strutil.Transform(inventory, func(item string) string {
				if ok, value := parseKevlar(item); ok {
					// Side effect of setting the armor value for the player
					p.State.Armor = value
					return EQUIPMENT_ID_KEVLAR
				}
				return item
			})

			// Flashbangs are only listed once when leaving the buyzone,
			// even if the player has 2 of them. In that case a second
			// one is added for clarity
			if strutil.Count(p.Equipment, EQUIPMENT_ID_FLASHBANG) == 2 {
				inventory = append(inventory, EQUIPMENT_ID_FLASHBANG)
			}

			// Sanity check to see that inventory is tracked correctly
			if !strutil.Equals(inventory, p.Equipment, func(item string) string {
				if ok, _ := parseKnife(item); ok {
					return EQUIPMENT_ID_KNIFE
				}
				return item
			}) {
				log.Printf("Player `%v` has equipment %+v but left buyzone with %+v\n", p.Nickname, p.Equipment, inventory)
			}

			// Set the player's inventory to be the buyzone summary
			p.Equipment = inventory
			return p
		})

	}
	s.transitPhases(_event)
}
