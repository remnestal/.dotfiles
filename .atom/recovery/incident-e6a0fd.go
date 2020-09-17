package incident

import (
	"fmt"

	"gitlab.com/abios/motattack/pkg/csgods/events"
)

// Incident defines a string ID representing a type of incident.
type Incident string

const (
	UnknownPlayer Incident = "unknown_player"
)

// Report represents a specific incident instance.
type Report struct {
	IncidentType Incident
	Origin       events.E
	Body         error
}

// Error provides a string representation of the incident instance.
func (r Report) Error() string {
	return fmt.Sprintf("%+v", r.Body)
}

// Description defines a string representing some text value providing
// additional information for understanding the incident.
type Description string

// Error provides a string representation of the incident description.
func (d Description) Error() string {
	return string(d)
}

// Player represents the player that is the subject of the incident.
type Player struct {
	Nickname, SteamId string
}

// Error provides a string representation of the player.
func (p Player) Error() string {
	return fmt.Sprintf("%+v", p)
}

// Expectation represents the failed expectation which caused the creation of
// the incident instance.
type Expectation struct {
	Expected interface{}
	Actual   interface{}
}

// Error provides a string representation of the failed expectation.
func (e Expectation) Error() string {
	return fmt.Sprintf("%+v", e)
}
