package state

import (
	"fmt"

	"gitlab.com/abios/motattack/pkg/csgods/events"
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
	Phase  string `json:"phase"`
	Length uint   `json:"length_minutes"`

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
		Phase string `json:"phase"`
	} `json:"round"`

	Home Team `json:"home"`
	Away Team `json:"away"`

	// Rounds holds the conlusion of all previosly played rounds
	Rounds []Round `json:"rounds"`
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

	Equipment []string `json:"equipment"`
}

// Commence resets all stats for all players and roundnumber etc and sets all phases to live
func (s *State) Commence() {
	s.Phase = "live"
	s.Length = 0

	s.Round.Number = 1
	s.Round.Phase = "live"

	s.Home.Score = 0
	for i := 0; i < len(s.Home.Players); i++ {
		s.Home.Players[i].Stats.Kills = 0
		s.Home.Players[i].Stats.Assists = 0
		s.Home.Players[i].Stats.Deaths = 0
		s.Home.Players[i].Stats.Adr = 0

		s.Home.Players[i].State.Health = 100
		s.Home.Players[i].State.Armor = 0
		s.Home.Players[i].State.Helmet = false
		s.Home.Players[i].State.Money = 800 // TODO: What is the starting money?

		s.Home.Players[i].State.Kills.Bodyshot = 0
		s.Home.Players[i].State.Kills.Headshot = 0

		s.Home.Players[i].State.Damage.Health = 0
		s.Home.Players[i].State.Damage.Armor = 0

		s.Home.Players[i].Equipment = []string{}
	}

	s.Away.Score = 0
	for i := 0; i < len(s.Away.Players); i++ {
		s.Away.Players[i].Stats.Kills = 0
		s.Away.Players[i].Stats.Assists = 0
		s.Away.Players[i].Stats.Deaths = 0
		s.Away.Players[i].Stats.Adr = 0

		s.Away.Players[i].State.Health = 100
		s.Away.Players[i].State.Armor = 0
		s.Away.Players[i].State.Helmet = false
		s.Away.Players[i].State.Money = 800 // TODO: What is the starting money?

		s.Away.Players[i].State.Kills.Bodyshot = 0
		s.Away.Players[i].State.Kills.Headshot = 0

		s.Away.Players[i].State.Damage.Health = 0
		s.Away.Players[i].State.Damage.Armor = 0

		s.Away.Players[i].Equipment = []string{}
	}

	s.Rounds = []Round{}
}

func NewState() *State {
	return &State{
		Home: Team{
			Faction: "ct",
		},
		Away: Team{
			Faction: "t",
		},
	}
}

func (s *State) UpdateSide(side string, transform func(Team) Team) {
	switch side {
	case s.Home.Faction:
		fmt.Println("updating ", s.Home.Faction)
		s.Home = transform(s.Home)
	case s.Away.Faction:
		fmt.Println("updating ", s.Away.Faction)
		s.Away = transform(s.Away)
	}
}

func (s *State) UpdatePlayer(player events.Player, transform func(Player) Player) {
	fmt.Printf("player %v on %v\n", player.SteamId, player.Side)
	s.UpdateSide(player.Side, func(current Team) Team {
		found := false
		for i, p := range current.Players {
			if player.SteamId == p.SteamId {
				fmt.Printf("updating existing player %v on %v\n", player.SteamId, player.Side)
				current.Players[i] = transform(p)
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("adding player %v to %v\n", player.SteamId, player.Side)
			current.Players = append(current.Players, transform(Player{
				Nickname: player.Nickname,
				SteamId:  player.SteamId,
			}))
		}
		return current
	})
}

func (s *State) Update(_event events.Event) {
	fmt.Printf("event %v\n", _event)
	switch event := _event.(type) {
	case *events.Cvar:
		// nothing
	case *events.GameOver:
		s.Map = event.Map
		s.Length = uint(event.Length)

		// TODO: Do something with score
		s.Phase = "gameover"
	case *events.Loading:
		s.Map = event.Map
	case *events.Molotov:
		// nothing
	case *events.PlayerAttacked:
		s.UpdatePlayer(event.Victim, func(current Player) Player {
			current.State.Health = event.Remaining.Health
			current.State.Armor = event.Remaining.Armor
			return current
		})

		s.UpdatePlayer(event.Player, func(current Player) Player {
			current.State.Damage.Health += event.Damage.Health
			current.State.Damage.Armor += event.Damage.Armor
			return current
		})
	case *events.PlayerKilledPlayer:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			current.Stats.Kills += 1
			if event.Headshot {
				current.State.Kills.Headshot += 1
			} else {
				current.State.Kills.Bodyshot += 1
			}
			return current
		})

		s.UpdatePlayer(event.Victim, func(current Player) Player {
			current.Stats.Deaths += 1
			return current
		})
	case *events.PlayerDestroyed:
	case *events.PlayerEntered:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			return Player{
				Nickname: event.Player.Nickname,
				SteamId:  event.Player.SteamId,
			}
		})
	case *events.PlayerConnected:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			return Player{
				Nickname: event.Player.Nickname,
				SteamId:  event.Player.SteamId,
			}
		})
	case *events.PlayerDisconnected:
		s.UpdateSide(event.Player.Side, func(current Team) Team {
			for i, p := range current.Players {
				if p.SteamId == event.Player.SteamId {
					current.Players = append(current.Players[:i], current.Players[i+1:]...)
					break
				}
			}
			return current
		})
	case *events.PlayerSwitched:
		var existing *Player
		s.UpdateSide(event.From, func(current Team) Team {
			for i, p := range current.Players {
				if p.SteamId == event.Player.SteamId {
					existing = &p
					current.Players = append(current.Players[:i], current.Players[i+1:]...)
					break
				}
			}
			return current
		})

		// It might be the case that a player moved from UNASSIGNED in which case we just add
		// the player as-is
		if existing != nil {
			s.UpdateSide(event.To, func(current Team) Team {
				current.Players = append(current.Players, *existing)
				return current
			})
		} else {
			s.UpdateSide(event.To, func(current Team) Team {
				current.Players = append(current.Players, Player{
					Nickname: event.Player.Nickname,
					SteamId:  event.Player.SteamId,
				})
				return current
			})
		}
	case *events.PlayerPickedUp:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			// TODO: Unless it's kevlar/assaultsuit etc then it should update the State.Armor
			current.Equipment = append(current.Equipment, event.What)
			return current
		})
	case *events.PlayerDropped:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			for i, e := range current.Equipment {
				if e == event.What {
					current.Equipment = append(current.Equipment[:i], current.Equipment[i+1:]...)
					break
				}
			}
			return current
		})
	case *events.PlayerLeftBuyzone:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			current.Equipment = event.With // TODO: Surely armor has to be handled
			return current
		})
	case *events.PlayerAquiredBomb:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			current.Equipment = append(current.Equipment, "c4")
			return current
		})
	case *events.PlayerDroppedBomb:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			for i, e := range current.Equipment {
				if e == "c4" {
					current.Equipment = append(current.Equipment[:i], current.Equipment[i+1:]...)
					break
				}
			}
			return current
		})
	case *events.PlayerPlantedBomb:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			for i, e := range current.Equipment {
				if e == "c4" {
					current.Equipment = append(current.Equipment[:i], current.Equipment[i+1:]...)
					break
				}
			}
			return current
		})
		s.Round.Phase = "bomb"
	case *events.PlayerDefusingBomb:
		s.Round.Phase = "defuse"
	case *events.PlayerDefusedBomb:
	case *events.PlayerKilledByBomb:
	case *events.PlayerSuicide:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			current.Stats.Deaths += 1
			return current
		})
	case *events.PlayerAssisted:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			current.Stats.Assists += 1
			return current
		})
	case *events.PlayerBlinded:
		// noop
	case *events.PlayerEconomyChange:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			current.State.Money = event.After
			if event.Delta > 0 {
				current.State.EquipValue += event.Delta
			}
			return current
		})
	case *events.PlayerThrew:
		s.UpdatePlayer(event.Player, func(current Player) Player {
			for i, e := range current.Equipment {
				if e == event.What {
					current.Equipment = append(current.Equipment[:i], current.Equipment[i+1:]...)
					break
				}
			}
			return current
		})
	case *events.StartedMap:
		s.Map = event.Map
		s.Commence()
		// TODO: This doesn't really mean it's live. Should probably check that there is 5
		// players and whatnot
		//s.Phase = "live"
	case *events.StartFreezetime:
		s.Round.Phase = "freezetime"
	case *events.TeamScored:
		s.UpdateSide(event.Team, func(current Team) Team {
			current.Score = event.Score
			return current
		})
	case *events.TeamWon:
		s.Rounds = append(s.Rounds, Round{
			Number: s.Round.Number,
			Winner: event.Team,
			Reason: event.Reason,
		})

		s.UpdateSide("t", func(current Team) Team {
			current.Score = event.Score.T
			return current
		})

		s.UpdateSide("ct", func(current Team) Team {
			current.Score = event.Score.CT
			return current
		})
	case *events.GameCommencing:
		s.Commence()
	case *events.MatchStarted:
		// TODO: When does this happen? What is this?
		s.Map = event.Map
		s.Commence()
	case *events.RoundStarted:
		if s.Phase == "live" {
			s.Round.Number += 1
		}
		s.Round.Phase = "live"
	case *events.RoundEnded:
		// TODO: Maybe increment round here? If it's live that is
		s.Round.Phase = "over"
	}
	fmt.Printf("%v\n", s.Home.Players)
	fmt.Printf("%v\n", s.Away.Players)
	fmt.Println("===============================")
}
