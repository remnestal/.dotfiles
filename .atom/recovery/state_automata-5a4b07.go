package state

import (
	"encoding/json"

	"gitlab.com/abios/motattack/pkg/csgods/events"
)

const (
	WARM_UP      string = "warmup"
	FREEZE_TIME  string = "freezetime"
	INTERMISSION string = "intermission"
	PRE_GAME     string = "pregame"
	OVER         string = "over"
	GAME_OVER    string = "gameover"
)

func lastActive(states []string) string {
	for i := len(states) - 1; i >= 0; i-- {
		switch states[i] {
		case "live", "warmup", "pregame": // warmup, pregame might not fit here
			return states[i]
		}
	}
	panic("No previous active state")
}

type Automata interface {
	MarshalJSON() ([]byte, error)
	Transist(_event events.Event)
	String() string
}

type past []string

type RoundAutomata struct {
	state       string
	transitions past
}

func (p *past) add(s string) {
	*p = append(*p, s)
}

func (auto *RoundAutomata) String() string {
	if auto.state != "" {
		return auto.state
	} else {
		return "undetermined"
	}
}

func (auto *RoundAutomata) MarshalJSON() ([]byte, error) {
	return json.Marshal(auto.String())
}

func (auto *RoundAutomata) Transist(_event events.Event) {
	switch _event.(type) {
	case *events.Loading, *events.LogStarted, *events.LogClosed:
		auto.transitions = []string{}
		auto.state = PRE_GAME

	case *events.StartedMap:
		auto.transitions.add(auto.state)
		auto.state = WARM_UP

	case *events.StartFreezetime:
		auto.transitions.add(auto.state)
		auto.state = FREEZE_TIME

	case *events.RoundStarted:
		if auto.state == FREEZE_TIME {
			auto.transitions.add(auto.state)
			auto.state = lastActive(auto.transitions)
		}

	case *events.RoundEnded:
		auto.transitions.add(auto.state)
		auto.state = OVER

	case *events.GameOver:
		auto.transitions.add(auto.state)
		auto.state = OVER
	}
}

type MatchAutomata struct {
	state       string
	transitions []string
}

func (auto *MatchAutomata) String() string {
	if auto.state != "" {
		return auto.state
	} else {
		return "undetermined"
	}
}

func (auto *MatchAutomata) MarshalJSON() ([]byte, error) {
	return json.Marshal(auto.String())
}

func (auto *MatchAutomata) Transist(_event events.Event) {
	switch _event.(type) {
	case *events.Loading, *events.LogStarted, *events.LogClosed:
		auto.transitions = []string{}
		auto.state = PRE_GAME

	case *events.StartedMap:
		auto.transitions.add(auto.state)
		auto.state = WARM_UP

	case *events.GameOver:
		auto.transitions.add(auto.state)
		auto.state = GAME_OVER

	case *events.MatchPaused:
		auto.transitions.add(auto.state)
		auto.state = INTERMISSION

	case *events.MatchResumed:
		auto.transitions.add(auto.state)
		auto.state = lastActive(auto.transitions)
	}
}
