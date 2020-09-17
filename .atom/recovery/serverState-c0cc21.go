package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"gitlab.com/abios/motattack/cmd/tydlig/token"
	"gitlab.com/abios/motattack/pkg/csgods/events"
	"gitlab.com/abios/motattack/pkg/csgods/events/state"
	"gitlab.com/abios/motattack/pkg/msg"
)

var errUnknownEvent error = errors.New("could not recognize event")

type serverState struct {
	Origin msg.Origin
	State  state.State
}

func (s *serverState) process(message msg.Msg) {
	quirky_event, err := events.UnmarshalJSON(message.Payload)
	if err != nil {
		log.Printf("unable to unmarshal event: %v\n", err)
		return
	}

	events := s.State.Dequirk(quirky_event)

	quirk_visualization := show_quirks && (len(events) != 1 || !reflect.DeepEqual(events[0], quirky_event))

	if quirk_visualization {
		if representation, err := s.digest(quirky_event); err != nil {
			log.Printf("unable to digest event for quirk visualization: %v", err)
			quirk_visualization = false
		} else {
			fmt.Println("<<<<<<< QUIRK")
			fmt.Println(representation)
			fmt.Println("=======")
		}
	}

	for _, event := range events {
		s.State.Update(event)
		if sequence, err := s.digest(event); err != nil {
			log.Printf("unable to digest event: %v", err)
		} else {
			fmt.Println(sequence.String())
		}
	}
	if quirk_visualization {
		fmt.Println(">>>>>>> d3qu12k3d")
	}
}

func (s *serverState) digest(e events.Event) (token.Sequence, error) {
	seq := token.NewSequence().
		Add(token.Host{Value: s.Origin.HostPort()}).
		Add(token.Timestamp{Value: e.Timestamp()}).
		Add(token.MatchPhase{Value: s.State.Phase.String()}).
		Add(token.RoundPhase{Value: s.State.Round.Phase.String()})

	switch event := e.(type) {
	case *events.LogClosed:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.LogStarted:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.Cvar:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()}).
			Add(token.Text{Value: event.Cvar}).
			Add(token.Text{Value: event.Value, Emph: true})

	case *events.GameOver:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()}).
			Add(token.Text{Value: "match played for"}).
			Add(token.Duration{Value: s.State.Timeline.End.Sub(*s.State.Timeline.Start)})

	case *events.Loading:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()}).
			Add(token.Map{Value: event.Map})

	case *events.Molotov:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerAttacked:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()}).
			Add(token.Nickname{Value: event.Victim.Nickname}).
			Add(token.Side{Value: event.Victim.Side})

	case *events.PlayerKilledPlayer:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerDestroyed:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerValidated:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerEntered:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerConnected:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerDisconnected:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerChangedName:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerSwitched:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerPickedUp:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerPurchased:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerDropped:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerLeftBuyzone:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerAquiredBomb:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerDroppedBomb:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerPlantedBomb:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerDefusingBomb:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerDefusedBomb:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerKilledByBomb:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerSuicide:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerAssisted:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerBlinded:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerEconomyChange:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerHas:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerThrew:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.PlayerChat:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.RestartRequired:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.CvarStart:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.CvarEnd:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.ServerMessage:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.StartedMap:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.FreezeTimeStarted:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.TeamScored:
		seq.Add(token.Nickname{Value: ""}).
			Add(token.Side{Value: event.Team}).
			Add(token.EventType{Value: event.Type()})

	case *events.TeamWon:
		seq.Add(token.Nickname{Value: ""}).
			Add(token.Side{Value: event.Team}).
			Add(token.EventType{Value: event.Type()})

	case *events.TeamPlaying:
		seq.Add(token.Nickname{Value: ""}).
			Add(token.Side{Value: event.Team}).
			Add(token.EventType{Value: event.Type()})

	case *events.GameCommencing:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.MatchStarted:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.MatchPaused:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.MatchResumed:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.MatchReloaded:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.RoundStarted:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.RoundEnded:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.RoundRestart:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.RoundDraw:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.EncryptionKeyEvent:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.VoteStarted:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.VoteSucceeded:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.VoteCast:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.VoteFailed:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.ErrorMessage:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.RconBadPassword:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.RconCommands:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.Accolade:
		seq.Add(token.Subject{token.Nickname{Value: event.Player.Nickname}}).
			Add(token.Side{Value: event.Player.Side}).
			Add(token.EventType{Value: event.Type()})

	case *events.BotMessage:
		seq.Add(token.Bot{}).
			Add(token.EventType{Value: event.Type()})

	case *events.LoadedPlugins:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	case *events.SteamAuth:
		seq.Add(token.World{}).
			Add(token.EventType{Value: event.Type()})

	default:
		return *seq, fmt.Errorf("%w: %v", errUnknownEvent, string(event.Type()))
	}
	return *seq, nil
}
