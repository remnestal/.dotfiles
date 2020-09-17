package main

import (
	"log"
	"time"

	"experiments/database/queries"
	"experiments/structs"
)

type PlayerBackground struct {
	Wins, Losses int64
}

func main() {
	log.Println("CS:GO Match Result Prediction")

	// Fetch all match results
	start := time.Now()
	matches, err := queries.GetMatchTeamResults()
	log.Printf("Fetched %v match results (%v)", len(matches), time.Since(start))
	if err != nil {
		log.Fatalf("Unable fetch match results from database: %v", err)
	}

	// What now? Somehow play the entire timeline of matches and for each match
	// make a prediction about who will win. This is something that should be
	// abstracted so that different heuristics for odds calculation can be applied
	// and compared.

	player_timeline := make(map[structs.PlayerId]map[structs.ReplayId]PlayerBackground)
	player_aggregate := make(map[structs.PlayerId]PlayerBackground)

	for _, m := range matches {
		if m.Winners == nil || m.Losers == nil {
			continue
		}
		for _, p := range *m.Winners {
			basis := player_aggregate[structs.PlayerId(p)]
			player_timeline, exist := player_timeline[structs.PlayerId(p)]
			if !exist {
				player_timeline = make(map[structs.ReplayId]PlayerBackground)
			}
			// Set the
			player_timeline[structs.ReplayId(m.ReplayId)] = basis
		}
	}

}
