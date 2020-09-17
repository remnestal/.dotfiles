package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/abios/motattack/pkg/csgods/events"
	"gitlab.com/abios/motattack/pkg/csgods/incident"
)

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
	//stub
}

func TestValidateLogStarted(t *testing.T) {
	//stub
}

func TestValidateCvar(t *testing.T) {
	//stub
}

func TestValidateGameOver(t *testing.T) {
	//stub
}

func TestValidateLoading(t *testing.T) {
	//stub
}

func TestValidateMolotov(t *testing.T) {
	//stub
}

func TestValidateRestartRequired(t *testing.T) {
	//stub
}

func TestValidateCvarStart(t *testing.T) {
	//stub
}

func TestValidateCvarEnd(t *testing.T) {
	//stub
}

func TestValidateServerMessage(t *testing.T) {
	//stub
}

func TestValidateStartedMap(t *testing.T) {
	//stub
}

func TestValidateFreezeTimeStarted(t *testing.T) {
	//stub
}

func TestValidateTeamScored(t *testing.T) {
	//stub
}

func TestValidateTeamWon(t *testing.T) {
	//stub
}

func TestValidateTeamPlaying(t *testing.T) {
	//stub
}

func TestValidateGameCommencing(t *testing.T) {
	//stub
}

func TestValidateMatchStarted(t *testing.T) {
	//stub
}

func TestValidateMatchPaused(t *testing.T) {
	//stub
}

func TestValidateMatchResumed(t *testing.T) {
	//stub
}

func TestValidateMatchReloaded(t *testing.T) {
	//stub
}

func TestValidateRoundStarted(t *testing.T) {
	//stub
}

func TestValidateRoundEnded(t *testing.T) {
	//stub
}

func TestValidateRoundRestart(t *testing.T) {
	//stub
}

func TestValidateRoundDraw(t *testing.T) {
	//stub
}

func TestValidateEncryptionKeyEvent(t *testing.T) {
	//stub
}

func TestValidateVoteStarted(t *testing.T) {
	//stub
}

func TestValidateVoteSucceeded(t *testing.T) {
	//stub
}

func TestValidateVoteCast(t *testing.T) {
	//stub
}

func TestValidateVoteFailed(t *testing.T) {
	//stub
}

func TestValidateErrorMessage(t *testing.T) {
	//stub
}

func TestValidateRconBadPassword(t *testing.T) {
	//stub
}

func TestValidateRconCommands(t *testing.T) {
	//stub
}

func TestValidateAccolade(t *testing.T) {
	//stub
}

func TestValidateBotMessage(t *testing.T) {
	//stub
}

func TestValidateLoadedPlugins(t *testing.T) {
	//stub
}

func TestValidateSteamAuth(t *testing.T) {
	//stub
}

func TestValidatePlayerConnected(t *testing.T) {
	//stub
}

func TestValidatePlayerValidated(t *testing.T) {
	//stub
}

func TestValidatePlayerEntered(t *testing.T) {
	//stub
}

func TestValidatePlayerDisconnected(t *testing.T) {
	//stub
}

func TestValidatePlayerSwitched(t *testing.T) {
	//stub
}

func TestValidatePlayerPickedUp(t *testing.T) {
	//stub
}

func TestValidatePlayerPurchased(t *testing.T) {
	//stub
}

func TestValidatePlayerDropped(t *testing.T) {
	//stub
}

func TestValidatePlayerThrew(t *testing.T) {
	//stub
}

func TestValidatePlayerEconomyChange(t *testing.T) {
	//stub
}

func TestValidatePlayerLeftBuyzone(t *testing.T) {
	//stub
}

func TestValidatePlayerPlantedBomb(t *testing.T) {
	//stub
}

func TestValidatePlayerDefusedBomb(t *testing.T) {
	//stub
}

func TestValidatePlayerAttacked(t *testing.T) {
	//stub
}

func TestValidatePlayerDestroyed(t *testing.T) {
	//stub
}

func TestValidatePlayerKilledByBomb(t *testing.T) {
	//stub
}

func TestValidatePlayerKilledPlayer(t *testing.T) {
	//stub
}

func TestValidatePlayerSuicide(t *testing.T) {
	//stub
}

func TestValidatePlayerHas(t *testing.T) {
	//stub
}

func TestValidatePlayerAssisted(t *testing.T) {
	//stub
}

func TestValidatePlayerChangedName(t *testing.T) {
	//stub
}
