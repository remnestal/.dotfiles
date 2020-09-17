package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/abios/motattack/pkg/csgods/events"
	"gitlab.com/abios/motattack/pkg/csgods/incident"
)

// number_incidents_of_type is an abstraction for counting the number of
// occurances of a specific kind of incident.
func number_incidents_of_type(ty incident.Incident, reports []incident.Report) (count int) {
	for i := range reports {
		if reports[i].IncidentType == ty {
			count++
		}
	}
	return
}

func Test_number_incidents_of_type(t *testing.T) {
	t.Run("no reports at all", func(t *testing.T) {
		assert.Equal(t, 0, number_incidents_of_type(incident.WrongEventType, []incident.Report{}))
	})
	t.Run("only matching reports", func(t *testing.T) {
		assert.Equal(t, 3, number_incidents_of_type(incident.WrongEventType, []incident.Report{
			incident.Report{IncidentType: incident.WrongEventType},
			incident.Report{IncidentType: incident.WrongEventType},
			incident.Report{IncidentType: incident.WrongEventType},
		}))
	})
	t.Run("no matching reports", func(t *testing.T) {
		assert.Equal(t, 0, number_incidents_of_type(incident.WrongEventType, []incident.Report{
			incident.Report{IncidentType: incident.UnknownPlayer},
			incident.Report{IncidentType: incident.UnknownPlayer},
			incident.Report{IncidentType: incident.UnknownPlayer},
		}))
	})
	t.Run("some matching reports", func(t *testing.T) {
		assert.Equal(t, 1, number_incidents_of_type(incident.WrongEventType, []incident.Report{
			incident.Report{IncidentType: incident.UnknownPlayer},
			incident.Report{IncidentType: incident.WrongEventType},
			incident.Report{IncidentType: incident.UnknownPlayer},
		}))
	})
}

// no_incidents is an abstraction for making sure that there's no incidents of
// the specified type in the passed slice of incidents.
func no_incidents(ty incident.Incident, reports []incident.Report) bool {
	return number_incidents_of_type(ty, reports) == 0
}

func Test_no_incidents(t *testing.T) {
	t.Run("no reports at all", func(t *testing.T) {
		assert.True(t, no_incidents(incident.WrongEventType, []incident.Report{}))
	})
	t.Run("only matching reports", func(t *testing.T) {
		assert.False(t, no_incidents(incident.WrongEventType, []incident.Report{
			incident.Report{IncidentType: incident.WrongEventType},
			incident.Report{IncidentType: incident.WrongEventType},
			incident.Report{IncidentType: incident.WrongEventType},
		}))
	})
	t.Run("no matching reports", func(t *testing.T) {
		assert.True(t, no_incidents(incident.WrongEventType, []incident.Report{
			incident.Report{IncidentType: incident.UnknownPlayer},
			incident.Report{IncidentType: incident.UnknownPlayer},
			incident.Report{IncidentType: incident.UnknownPlayer},
		}))
	})
	t.Run("some matching reports", func(t *testing.T) {
		assert.False(t, no_incidents(incident.WrongEventType, []incident.Report{
			incident.Report{IncidentType: incident.UnknownPlayer},
			incident.Report{IncidentType: incident.WrongEventType},
			incident.Report{IncidentType: incident.UnknownPlayer},
		}))
	})
}

// one_incident is an abstraction for making sure that there's only one incident
// of the specified type in the passed slice of incidents.
func one_incident(ty incident.Incident, reports []incident.Report) bool {
	return number_incidents_of_type(ty, reports) == 1
}

func Test_one_incident(t *testing.T) {
	t.Run("no reports at all", func(t *testing.T) {
		assert.False(t, one_incident(incident.WrongEventType, []incident.Report{}))
	})
	t.Run("only matching reports", func(t *testing.T) {
		assert.False(t, one_incident(incident.WrongEventType, []incident.Report{
			incident.Report{IncidentType: incident.WrongEventType},
			incident.Report{IncidentType: incident.WrongEventType},
			incident.Report{IncidentType: incident.WrongEventType},
		}))
	})
	t.Run("no matching reports", func(t *testing.T) {
		assert.False(t, one_incident(incident.WrongEventType, []incident.Report{
			incident.Report{IncidentType: incident.UnknownPlayer},
			incident.Report{IncidentType: incident.UnknownPlayer},
			incident.Report{IncidentType: incident.UnknownPlayer},
		}))
	})
	t.Run("some matching reports", func(t *testing.T) {
		assert.False(t, one_incident(incident.WrongEventType, []incident.Report{
			incident.Report{IncidentType: incident.UnknownPlayer},
			incident.Report{IncidentType: incident.WrongEventType},
			incident.Report{IncidentType: incident.WrongEventType},
			incident.Report{IncidentType: incident.UnknownPlayer},
		}))
	})
	t.Run("one matching reports", func(t *testing.T) {
		assert.True(t, one_incident(incident.WrongEventType, []incident.Report{
			incident.Report{IncidentType: incident.UnknownPlayer},
			incident.Report{IncidentType: incident.WrongEventType},
			incident.Report{IncidentType: incident.UnknownPlayer},
		}))
	})
}

func Test_validate_event_type(t *testing.T) {
	s := testState()
	t.Run("correct event type yields no incident", func(t *testing.T) {
		assert.Equal(t, []incident.Report{}, s.validate_event_type(&events.PlayerPickedUp{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerPickedUpEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
			What: "kevlar",
		}, events.PlayerPickedUpEventType))
	})
	t.Run("unknown player yields an incident", func(t *testing.T) {
		assert.Equal(t, []incident.Report{
			incident.Report{
				IncidentType: incident.WrongEventType,
				Origin:       nil,
				Body: incident.Expectation{
					Expected: events.PlayerPickedUpEventType,
					Actual:   events.PlayerDroppedEventType,
				},
			},
		}, s.validate_event_type(&events.PlayerPickedUp{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerDroppedEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
			What: "kevlar",
		}, events.PlayerPickedUpEventType))
	})
}

func Test_validate_player_existance(t *testing.T) {
	s := testState()
	t.Run("known player yields no incident", func(t *testing.T) {
		assert.Equal(t, []incident.Report{}, s.validate_player_existance(events.Player{
			SteamId: "player_0",
		}))
	})
	t.Run("unknown player yields an incident", func(t *testing.T) {
		assert.Equal(t, []incident.Report{
			incident.Report{
				IncidentType: incident.UnknownPlayer,
				Origin:       nil,
				Subject: incident.Player{
					SteamId:  "player_?",
					Nickname: "pupdog",
					Side:     "ct",
				},
				Body: incident.Description("unknown player unexpectedly listed in event"),
			},
		}, s.validate_player_existance(events.Player{
			SteamId:  "player_?",
			Nickname: "pupdog",
			Side:     "ct",
		}))
	})
}

func Test_validate_player_nickname(t *testing.T) {
	s := testState()
	t.Run("known player name yields no incident", func(t *testing.T) {
		assert.Equal(t, []incident.Report{}, s.validate_player_nickname(events.Player{
			SteamId:  "player_0",
			Nickname: "p0",
		}))
	})
	t.Run("unreported player name yields an incident", func(t *testing.T) {
		assert.Equal(t, []incident.Report{
			incident.Report{
				IncidentType: incident.UnreportedPlayerNameChange,
				Origin:       nil,
				Subject: incident.Player{
					SteamId:  "player_0",
					Nickname: "pupdog",
					Side:     "ct",
				},
				Body: incident.Expectation{
					Expected: "p0",
					Actual:   "pupdog",
				},
			},
		}, s.validate_player_nickname(events.Player{
			SteamId:  "player_0",
			Nickname: "pupdog",
			Side:     "ct",
		}))
	})
}

func Test_validate_player_faction(t *testing.T) {
	s := testState()
	t.Run("known player faction yields no incident", func(t *testing.T) {
		assert.Equal(t, []incident.Report{}, s.validate_player_faction(events.Player{
			SteamId: "player_0",
			Side:    "ct",
		}))
	})
	t.Run("unreported player faction yields an incident", func(t *testing.T) {
		assert.Equal(t, []incident.Report{
			incident.Report{
				IncidentType: incident.UnreportedPlayerSideSwitch,
				Origin:       nil,
				Subject: incident.Player{
					SteamId:  "player_0",
					Nickname: "p0",
					Side:     "t",
				},
				Body: incident.Expectation{
					Expected: "ct",
					Actual:   "t",
				},
			},
		}, s.validate_player_faction(events.Player{
			SteamId:  "player_0",
			Nickname: "p0",
			Side:     "t",
		}))
	})
}

func Test_flatten(t *testing.T) {
	t.Run("slices are combined", func(t *testing.T) {
		slice1 := []incident.Report{
			incident.Report{IncidentType: incident.UnknownPlayer},
		}
		slice2 := []incident.Report{
			incident.Report{IncidentType: incident.UnreportedPlayerNameChange},
			incident.Report{IncidentType: incident.UnreportedPlayerSideSwitch},
		}
		assert.Equal(t, []incident.Report{
			incident.Report{IncidentType: incident.UnknownPlayer},
			incident.Report{IncidentType: incident.UnreportedPlayerNameChange},
			incident.Report{IncidentType: incident.UnreportedPlayerSideSwitch},
		}, flatten(slice1, slice2))
	})
}

func TestValidateLogClosed(t *testing.T) {
	base_event := func() events.LogClosed {
		return events.LogClosed{
			E: events.E{
				Ty: events.LogClosedEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.LogClosedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateLogStarted(t *testing.T) {
	base_event := func() events.LogStarted {
		return events.LogStarted{
			E: events.E{
				Ty: events.LogStartedEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.LogStartedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateCvar(t *testing.T) {
	base_event := func() events.Cvar {
		return events.Cvar{
			E: events.E{
				Ty: events.CvarEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.CvarEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateGameOver(t *testing.T) {
	base_event := func() events.GameOver {
		return events.GameOver{
			E: events.E{
				Ty: events.GameOverEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		},
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.GameOverEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateLoading(t *testing.T) {
	base_event := func() events.Loading {
		return events.Loading{
			E: events.E{
				Ty: events.LoadingEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.LoadingEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateMolotov(t *testing.T) {
	base_event := func() events.Molotov {
		return events.Molotov{
			E: events.E{
				Ty: events.MolotovEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.MolotovEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateRestartRequired(t *testing.T) {
	base_event := func() events.RestartRequired {
		return events.RestartRequired{
			E: events.E{
				Ty: events.RestartRequiredEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.RestartRequiredEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateCvarStart(t *testing.T) {
	base_event := func() events.CvarStart {
		return events.CvarStart{
			E: events.E{
				Ty: events.CvarStartEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.CvarStartEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateCvarEnd(t *testing.T) {
	base_event := func() events.CvarEnd {
		return events.CvarEnd{
			E: events.E{
				Ty: events.CvarEndEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.CvarEndEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateServerMessage(t *testing.T) {
	base_event := func() events.ServerMessage {
		return events.ServerMessage{
			E: events.E{
				Ty: events.ServerMessageEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.ServerMessageEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateStartedMap(t *testing.T) {
	base_event := func() events.StartedMap {
		return events.StartedMap{
			E: events.E{
				Ty: events.StartedMapEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.StartedMapEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateFreezeTimeStarted(t *testing.T) {
	base_event := func() events.FreezeTimeStarted {
		return events.FreezeTimeStarted{
			E: events.E{
				Ty: events.FreezeTimeStartedEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.FreezeTimeStartedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateTeamScored(t *testing.T) {
	base_event := func() events.TeamScored {
		return events.TeamScored{
			TeamEvent: events.TeamEvent{
				E: events.E{
					Ty: events.TeamScoredEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.TeamScoredEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateTeamWon(t *testing.T) {
	base_event := func() events.TeamWon {
		return events.TeamWon{
			TeamEvent: events.TeamEvent{
				E: events.E{
					Ty: events.TeamWonEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.TeamWonEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateTeamPlaying(t *testing.T) {
	base_event := func() events.TeamPlaying {
		return events.TeamPlaying{
			TeamEvent: events.TeamEvent{
				E: events.E{
					Ty: events.TeamPlayingEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.TeamPlayingEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateGameCommencing(t *testing.T) {
	base_event := func() events.GameCommencing {
		return events.GameCommencing{
			E: events.E{
				Ty: events.GameCommencingEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.GameCommencingEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateMatchStarted(t *testing.T) {
	base_event := func() events.MatchStarted {
		return events.MatchStarted{
			E: events.E{
				Ty: events.MatchStartedEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.MatchStartedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateMatchPaused(t *testing.T) {
	base_event := func() events.MatchPaused {
		return events.MatchPaused{
			E: events.E{
				Ty: events.MatchPausedEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.MatchPausedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateMatchResumed(t *testing.T) {
	base_event := func() events.MatchResumed {
		return events.MatchResumed{
			E: events.E{
				Ty: events.MatchResumedEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.MatchResumedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateMatchReloaded(t *testing.T) {
	base_event := func() events.MatchReloaded {
		return events.MatchReloaded{
			E: events.E{
				Ty: events.MatchReloadedEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.MatchReloadedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateRoundStarted(t *testing.T) {
	base_event := func() events.RoundStarted {
		return events.RoundStarted{
			E: events.E{
				Ty: events.RoundStartedEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.RoundStartedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateRoundEnded(t *testing.T) {
	base_event := func() events.RoundEnded {
		return events.RoundEnded{
			E: events.E{
				Ty: events.RoundEndedEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.RoundEndedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateRoundRestart(t *testing.T) {
	base_event := func() events.RoundRestart {
		return events.RoundRestart{
			E: events.E{
				Ty: events.RoundRestartEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.RoundRestartEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateRoundDraw(t *testing.T) {
	base_event := func() events.RoundDraw {
		return events.RoundDraw{
			E: events.E{
				Ty: events.RoundDrawEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.RoundDrawEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateEncryptionKeyEvent(t *testing.T) {
	base_event := func() events.EncryptionKeyEvent {
		return events.EncryptionKeyEvent{
			E: events.E{
				Ty: events.EncryptionKeyEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.EncryptionKeyEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateVoteStarted(t *testing.T) {
	base_event := func() events.VoteStarted {
		return events.VoteStarted{
			VoteEvent: events.VoteEvent{
				E: events.E{
					Ty: events.VoteStartedEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.VoteStartedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateVoteSucceeded(t *testing.T) {
	base_event := func() events.VoteSucceeded {
		return events.VoteSucceeded{
			VoteEvent: events.VoteEvent{
				E: events.E{
					Ty: events.VoteSucceededEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.VoteSucceededEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateVoteCast(t *testing.T) {
	base_event := func() events.VoteCast {
		return events.VoteCast{
			VoteEvent: events.VoteEvent{
				E: events.E{
					Ty: events.VoteCastEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.VoteCastEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateVoteFailed(t *testing.T) {
	base_event := func() events.VoteFailed {
		return events.VoteFailed{
			VoteEvent: events.VoteEvent{
				E: events.E{
					Ty: events.VoteFailedEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.VoteFailedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateErrorMessage(t *testing.T) {
	base_event := func() events.ErrorMessage {
		return events.ErrorMessage{
			E: events.E{
				Ty: events.ErrorMessageEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.ErrorMessageEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateRconBadPassword(t *testing.T) {
	base_event := func() events.RconBadPassword {
		return events.RconBadPassword{
			RconEvent: events.RconEvent{
				E: events.E{
					Ty: events.RconBadPasswordEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.RconBadPasswordEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateRconCommands(t *testing.T) {
	base_event := func() events.RconCommands {
		return events.RconCommands{
			RconEvent: events.RconEvent{
				E: events.E{
					Ty: events.RconCommandsEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.RconCommandsEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateAccolade(t *testing.T) {
	base_event := func() events.Accolade {
		return events.Accolade{
			E: events.E{
				Ty: events.AccoladeEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.AccoladeEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateBotMessage(t *testing.T) {
	base_event := func() events.BotMessage {
		return events.BotMessage{
			E: events.E{
				Ty: events.BotMessageEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.BotMessageEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateLoadedPlugins(t *testing.T) {
	base_event := func() events.LoadedPlugins {
		return events.LoadedPlugins{
			E: events.E{
				Ty: events.LoadedPluginsEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.LoadedPluginsEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidateSteamAuth(t *testing.T) {
	base_event := func() events.SteamAuth {
		return events.SteamAuth{
			E: events.E{
				Ty: events.SteamAuthEventType,
				Ti: time.Time{}.Add(123 * time.Second),
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.SteamAuthEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerConnected(t *testing.T) {
	base_event := func() events.PlayerConnected {
		return events.PlayerConnected{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerConnectedEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerConnectedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerValidated(t *testing.T) {
	base_event := func() events.PlayerValidated {
		return events.PlayerValidated{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerValidatedIdEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerValidatedIdEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerEntered(t *testing.T) {
	base_event := func() events.PlayerEntered {
		return events.PlayerEntered{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerEnteredGameEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerEnteredGameEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerDisconnected(t *testing.T) {
	base_event := func() events.PlayerDisconnected {
		return events.PlayerDisconnected{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerDisconnectedEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerDisconnectedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerSwitched(t *testing.T) {
	base_event := func() events.PlayerSwitched {
		return events.PlayerSwitched{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerSwitchedEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerSwitchedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerPickedUp(t *testing.T) {
	base_event := func() events.PlayerPickedUp {
		return events.PlayerPickedUp{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerPickedUpEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerPickedUpEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerPurchased(t *testing.T) {
	base_event := func() events.PlayerPurchased {
		return events.PlayerPurchased{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerPurchasedEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerPurchasedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerDropped(t *testing.T) {
	base_event := func() events.PlayerDropped {
		return events.PlayerDropped{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerDroppedEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerDroppedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerThrew(t *testing.T) {
	base_event := func() events.PlayerThrew {
		return events.PlayerThrew{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerThrewEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerThrewEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerEconomyChange(t *testing.T) {
	base_event := func() events.PlayerEconomyChange {
		return events.PlayerEconomyChange{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerEconomyChangeEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerEconomyChangeEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerLeftBuyzone(t *testing.T) {
	base_event := func() events.PlayerLeftBuyzone {
		return events.PlayerLeftBuyzone{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerLeftBuyzoneEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerLeftBuyzoneEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerPlantedBomb(t *testing.T) {
	base_event := func() events.PlayerPlantedBomb {
		return events.PlayerPlantedBomb{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerPlantedBombEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerPlantedBombEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerDefusedBomb(t *testing.T) {
	base_event := func() events.PlayerDefusedBomb {
		return events.PlayerDefusedBomb{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerDefusedBombEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerDefusedBombEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerAttacked(t *testing.T) {
	base_event := func() events.PlayerAttacked {
		return events.PlayerAttacked{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerAttackedEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerAttackedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerDestroyed(t *testing.T) {
	base_event := func() events.PlayerDestroyed {
		return events.PlayerDestroyed{
			PlayerKilled: events.PlayerKilled{
				PlayerEvent: events.PlayerEvent{
					E: events.E{
						Ty: events.PlayerDestroyedEventType,
						Ti: time.Time{}.Add(123 * time.Second),
					},
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerDestroyedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerKilledByBomb(t *testing.T) {
	base_event := func() events.PlayerKilledByBomb {
		return events.PlayerKilledByBomb{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerKilledByBombEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerKilledByBombEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerKilledPlayer(t *testing.T) {
	base_event := func() events.PlayerKilledPlayer {
		return events.PlayerKilledPlayer{
			PlayerKilled: events.PlayerKilled{
				PlayerEvent: events.PlayerEvent{
					E: events.E{
						Ty: events.PlayerKilledPlayerEventType,
						Ti: time.Time{}.Add(123 * time.Second),
					},
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerKilledPlayerEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerSuicide(t *testing.T) {
	base_event := func() events.PlayerSuicide {
		return events.PlayerSuicide{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerSuicideEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerSuicideEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerHas(t *testing.T) {
	base_event := func() events.PlayerHas {
		return events.PlayerHas{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerHasEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerHasEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerAssisted(t *testing.T) {
	base_event := func() events.PlayerAssisted {
		return events.PlayerAssisted{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerAssistedEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerAssistedEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}

func TestValidatePlayerChangedName(t *testing.T) {
	base_event := func() events.PlayerChangedName {
		return events.PlayerChangedName{
			PlayerEvent: events.PlayerEvent{
				E: events.E{
					Ty: events.PlayerChangedNameEventType,
					Ti: time.Time{}.Add(123 * time.Second),
				},
			},
		}
	}
	t.Run("event type", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			s := New()
			event := base_event()
			assert.True(t, no_incidents(incident.WrongEventType, s.Validate(&event))
		})
		t.Run("wrong", func(t *testing.T) {
			s := New()
			event := base_event()
			event.E.Ty = events.EventType("bad event type")
			assert.Equal(t, []incident.Report{
				incident.Report{
					IncidentType: incident.WrongEventType,
					Origin:       &event,
					Body: incident.Expectation{
						Expected: events.PlayerChangedNameEventType,
						Actual:   events.EventType("bad event type"),
					},
				},
			}, s.Validate(&event))
		})
	})
}
