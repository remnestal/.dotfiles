package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"experiments/database/queries"
	"experiments/fileio"
	"experiments/hltv"
	"experiments/structs"
	"experiments/structs/svm"

	hltvpkg "github.com/remnestal/hltv"
)

const WEIGHTS_SUFFIX string = ".weights"

var (
	partitions_dir  string
	identities_file string
	output_file     string
)

func init() {
	flag.StringVar(&partitions_dir, "weights", "/tmp/weights", "directory with weights")
	flag.StringVar(&identities_file, "ids", "/tmp/identities", "file with HLTV identities")
	flag.StringVar(&output_file, "output", "/tmp/evaluation.hltv.csv", "output CSV file")
	flag.Parse()
}

type Partition struct {
	Start, End time.Time
	Weights    map[string]float64
}

func load_partitions() ([]Partition, error) {
	files, err := ioutil.ReadDir(partitions_dir)
	if err != nil {
		return nil, fmt.Errorf("Unable to read directory %q: %w", partitions_dir, err)
	}
	partitions := make([]Partition, 0)
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, WEIGHTS_SUFFIX) {
			partition := Partition{}

			// Parse start- and end-date from file name
			components := strings.Split(strings.TrimSuffix(name, WEIGHTS_SUFFIX), "-")
			if partition.Start, err = time.Parse("2006-01-02", strings.Join(components[2:5], "-")); err != nil {
				return nil, fmt.Errorf("Unable to parse start-date of weight file %q: %w", name, err)
			}
			if partition.End, err = time.Parse("2006-01-02", strings.Join(components[5:8], "-")); err != nil {
				return nil, fmt.Errorf("Unable to parse end-date of weight file %q: %w", name, err)
			}

			// Read and parse the weights from file
			partition.Weights, err = fileio.ReadWeights(fmt.Sprintf("%v/%v", partitions_dir, name))
			if err != nil {
				log.Fatalf("Unable to read weight vector: %v", err)
			}

			partitions = append(partitions, partition)
		}
	}
	return partitions, nil
}

type feature_query func(start time.Time, end time.Time, playerIds ...structs.PlayerId) (svm.PlayerAggregates, error)

func db_query(desc string, start, end time.Time, query feature_query, ids ...structs.PlayerId) svm.PlayerAggregates {
	timer_start := time.Now()
	if aggregate, err := query(start, end, ids...); err != nil {
		log.Printf("[NON TERMINATING ERROR] Unable to fetch %q because: %v", desc, err)
		return svm.PlayerAggregates{}
	} else {
		log.Printf("aggregated %q for %v players in %v", desc, len(ids), time.Since(timer_start))
		return aggregate
	}
}

func NormalizedRatingsByPlayerId(stats svm.PlayerAggregates, weights map[string]float64) map[structs.PlayerId]float64 {
	ratings := make(map[structs.PlayerId]float64, len(stats))
	// Find the minimum and maximum rating
	var min, max float64
	for id, features := range stats {
		r := features.Rate(weights, nil)
		if r < min {
			min = r
		}
		if r > max {
			max = r
		}
		ratings[id] = r
	}
	// Normalize the players' ratings
	for id, rating := range ratings {
		ratings[id] = (rating - min) / (max - min)
	}
	return ratings
}

func main() {
	log.Println("HLTV Player Rating Comparison")

	// Load weight partitions from the specified directory
	partitions, err := load_partitions()
	if err != nil {
		log.Fatalf("Unable to load weights: %v", err)
	}

	log.Printf("Loaded %v weight partitions\n", len(partitions))

	// Load existing HLTV identities from file
	all_identities, err := hltv.ReadIdentitiesFile(identities_file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("Creating new identities file %q", identities_file)
			all_identities = structs.IdentityMap{}
			if err := hltv.WriteIdentitiesFile(identities_file, all_identities); err != nil {
				log.Fatalf("Unable to store identities: %v", err)
			}
		} else {
			log.Fatalf("Unable to load identities: %v", err)
		}
	}

	log.Printf("Loaded %v HLTV identities\n", len(all_identities))

	// Store the aggregated evaluation results for all of the partitions
	results := make([]Result, 0)
	timer := time.Now()

	// Process each partition separately
	for index, partition := range partitions {
		partition_timer := time.Now()

		log.Printf("Processing partition %v/%v ...", index+1, len(partitions))

		// Fetch HLTV stats for the given time period
		hltv_timer := time.Now()
		hltv_stats, err := hltvpkg.PlayerStats(hltvpkg.PlayerStatsParams{
			RequestConfig: hltvpkg.RequestConfig{
				Delay:    time.Second * 3,
				Scheme:   "https",
				Domain:   "www.hltv.org",
				BasePath: "/stats/players",
				BaseParams: fmt.Sprintf("startDate=%v&endDate=%v&rankingFilter=All",
					partition.Start.Format("2006-01-02"),
					partition.End.Format("2006-01-02")),
			},
			Verbose: false,
		}, hltvpkg.IdentityMap(all_identities))
		if err != nil {
			log.Fatalf("Unable to fetch HLTV player stats: %v", err)
		}

		// Update the HLTV identities file
		priori := len(all_identities)
		relevant_identities := make(structs.IdentityMap, len(hltv_stats.Players))
		for _, p := range hltv_stats.Players {
			// Store both in the "global" identity cache and make a separate one for
			// players relevant for this partition
			all_identities[hltvpkg.PlayerId(p.Id)] = p.PlayerNames
			relevant_identities[hltvpkg.PlayerId(p.Id)] = p.PlayerNames
		}
		if err := hltv.WriteIdentitiesFile(identities_file, all_identities); err != nil {
			log.Fatalf("Unable to store identities: %v", err)
		}

		log.Printf("Fetched %v player stats (%v new identities added, %v elapsed)", len(hltv_stats.Players), len(all_identities)-priori, time.Since(hltv_timer))

		// Attempt to pair HLTV identities
		paired, unknown, ambigious, err := hltv.Pair(relevant_identities)
		if err != nil {
			log.Fatalf("Unable to pair HLTV identities: %v", err)
		}

		log.Printf("Successfully paired: %v (%.f%%) HLTV players (%v unknown, %v ambiguous)",
			len(paired), (float64(len(paired))/float64(len(relevant_identities)))*100,
			len(unknown), len(ambigious))

		// Make a slice of the abios IDs for which features should be aggregated
		abios_ids := make([]structs.PlayerId, 0)
		for _, abios_id := range paired {
			abios_ids = append(abios_ids, abios_id)
		}

		// Fetch aggregated stats for each player
		features := db_query("damage dealt", partition.Start, partition.End, queries.DamageDealtStatsByPlayer, abios_ids...)
		features = features.Combine(db_query("damage taken", partition.Start, partition.End, queries.DamageTakenStatsByPlayer, abios_ids...))
		features = features.Combine(db_query("assists", partition.Start, partition.End, queries.AssistStatsByPlayer, abios_ids...))
		features = features.Combine(db_query("weapons", partition.Start, partition.End, queries.WeaponClassStatsByPlayer, abios_ids...))
		features = features.Combine(db_query("reloads", partition.Start, partition.End, queries.ReloadStatsByPlayer, abios_ids...))
		features = features.Combine(db_query("grenades", partition.Start, partition.End, queries.GrenadeStatsByPlayer, abios_ids...))

		// Calculate a normalized rating for each player in each domain
		ratings_abios := NormalizedRatingsByPlayerId(features, partition.Weights)

		evaluate := func(system string, getter hltv.RatingSelector) {
			analyse := func(domain string, evaluator hltv.AnalysisFunction) {
				ratings_hltv := hltv.NormalizedRatingsByPlayerId(paired, hltv_stats, getter)
				results = append(results, Result{
					PartitionIndex: int64(index),
					Start:          partition.Start.Unix(),
					End:            partition.End.Unix(),
					Domain:         domain,
					System:         system,
					Evaluation:     evaluator(ratings_hltv, ratings_abios),
				})
			}
			analyse("continuous", hltv.RegressionAnalysis)
			analyse("discreet", hltv.DiscreetAnalysis)
		}

		evaluate("rating", func(p hltvpkg.PlayerRow) float64 {
			return p.Rating
		})
		evaluate("kd_ratio", func(p hltvpkg.PlayerRow) float64 {
			return p.KDRatio
		})
		evaluate("kd_diff", func(p hltvpkg.PlayerRow) float64 {
			return float64(p.KDDiff)
		})

		log.Printf("Partition evaluated (%v elapsed)", time.Since(partition_timer))

	}

	log.Printf("All partitions evaluated (%v elapsed)", time.Since(timer))

	// Write results to file
	bytes := make([]byte, 0)
	bytes = append(bytes, []byte(fmt.Sprintf("%v\n", strings.Join(ResultHeaders(), ",")))...)
	for _, res := range results {
		bytes = append(bytes, []byte(fmt.Sprintf("%v\n", strings.Join(res.CSV(), ",")))...)
	}
	if err := ioutil.WriteFile(output_file, bytes, 0644); err != nil {
		log.Fatalf("Unable to write output file: %v", err)
	}

}
