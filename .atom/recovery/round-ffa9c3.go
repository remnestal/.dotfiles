package automata

import (
	"encoding/json"

	"gitlab.com/abios/motattack/pkg/csgods/events"
)

type RoundAutomata struct {
	state       string
	transitions past
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
		auto.transitions.Record(auto.state)
		auto.state = WARM_UP

	case *events.StartFreezetime:
		auto.transitions.Record(auto.state)
		auto.state = FREEZE_TIME

	case *events.RoundStarted:
		if auto.state == FREEZE_TIME {
			auto.transitions.Record(auto.state)
			auto.state = auto.transitions.LastRunning()
		}

	case *events.RoundEnded:
		auto.transitions.Record(auto.state)
		auto.state = OVER

	case *events.GameOver:
		auto.transitions.Record(auto.state)
		auto.state = OVER
	}
}
