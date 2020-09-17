package perplayer

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"experiments/database"
	"experiments/structs"
	"experiments/structs/svm"
)

type MatchPlayerFeatures map[structs.Eid]map[structs.Eid]svm.PlayerFeatures

func DamageGiven(matches, players []structs.Eid) (MatchPlayerFeatures, error) {
	features := make(MatchPlayerFeatures)
	stmt := `
  select
    CsRound.replay_id,
    CsDamage.given_eid,

    @dist = ST_Distance(Point(giver_pos_x,giver_pos_y), Point(taker_pos_x,taker_pos_x)) as distance,

    -- 500
    sum(if(hitgroup=1 and @dist > 0 and @dist <= 500, CsDamage.dmg_health, 0)) as d500_head_health,
    sum(if(not hitgroup=1 and @dist > 0 and @dist <= 500, CsDamage.dmg_health, 0)) as d500_body_health,
    sum(if(hitgroup=1 and @dist > 0 and @dist <= 500, CsDamage.dmg_armor, 0)) as d500_head_armor,
    sum(if(not hitgroup=1 and @dist > 0 and @dist <= 500, CsDamage.dmg_armor, 0)) as d500_body_armor,
    -- 1000
    sum(if(hitgroup=1 and @dist > 500 and @dist <= 1000, CsDamage.dmg_health, 0)) as d1000_head_health,
    sum(if(not hitgroup=1 and @dist > 500 and @dist <= 1000, CsDamage.dmg_health, 0)) as d1000_bod_healthy,
    sum(if(hitgroup=1 and @dist > 500 and @dist <= 1000, CsDamage.dmg_armor, 0)) as d1000_head_armor,
    sum(if(not hitgroup=1 and @dist > 500 and @dist <= 1000, CsDamage.dmg_armor, 0)) as d1000_body_armor,
    -- 1500
    sum(if(hitgroup=1 and @dist > 1000 and @dist <= 1500, CsDamage.dmg_health, 0)) as d1500_head_health,
    sum(if(not hitgroup=1 and @dist > 1000 and @dist <= 1500, CsDamage.dmg_health, 0)) as d1500_body_health,
    sum(if(hitgroup=1 and @dist > 1000 and @dist <= 1500, CsDamage.dmg_armor, 0)) as d1500_head_armor,
    sum(if(not hitgroup=1 and @dist > 1000 and @dist <= 1500, CsDamage.dmg_armor, 0)) as d1500_body_armor,
    -- 2000
    sum(if(hitgroup=1 and @dist > 1500 and @dist <= 2000, CsDamage.dmg_health, 0)) as d2000_head_health,
    sum(if(not hitgroup=1 and @dist > 1500 and @dist <= 2000, CsDamage.dmg_health, 0)) as d2000_body_health,
    sum(if(hitgroup=1 and @dist > 1500 and @dist <= 2000, CsDamage.dmg_armor, 0)) as d2000_head_armor,
    sum(if(not hitgroup=1 and @dist > 1500 and @dist <= 2000, CsDamage.dmg_armor, 0)) as d2000_body_armor,
    -- 2500
    sum(if(hitgroup=1 and @dist > 2000 and @dist <= 2500, CsDamage.dmg_health, 0)) as d2500_head_health,
    sum(if(not hitgroup=1 and @dist > 2000 and @dist <= 2500, CsDamage.dmg_health, 0)) as d2500_body_health,
    sum(if(hitgroup=1 and @dist > 2000 and @dist <= 2500, CsDamage.dmg_armor, 0)) as d2500_head_armor,
    sum(if(not hitgroup=1 and @dist > 2000 and @dist <= 2500, CsDamage.dmg_armor, 0)) as d2500_body_armor,
    -- 3000
    sum(if(hitgroup=1 and @dist > 2500 and @dist <= 3000, CsDamage.dmg_health, 0)) as d3000_head_health,
    sum(if(not hitgroup=1 and @dist > 2500 and @dist <= 3000, CsDamage.dmg_health, 0)) as d3000_body_health,
    sum(if(hitgroup=1 and @dist > 2500 and @dist <= 3000, CsDamage.dmg_armor, 0)) as d3000_head_armor,
    sum(if(not hitgroup=1 and @dist > 2500 and @dist <= 3000, CsDamage.dmg_armor, 0)) as d3000_body_armor,
    -- 3500
    sum(if(hitgroup=1 and @dist > 3000 and @dist <= 3500, CsDamage.dmg_health, 0)) as d3500_head_health,
    sum(if(not hitgroup=1 and @dist > 3000 and @dist <= 3500, CsDamage.dmg_health, 0)) as d3500_body_health,
    sum(if(hitgroup=1 and @dist > 3000 and @dist <= 3500, CsDamage.dmg_armor, 0)) as d3500_head_armor,
    sum(if(not hitgroup=1 and @dist > 3000 and @dist <= 3500, CsDamage.dmg_armor, 0)) as d3500_body_armor,
    -- 4000
    sum(if(hitgroup=1 and @dist > 3500 and @dist <= 4000, CsDamage.dmg_health, 0)) as d4000_head_health,
    sum(if(not hitgroup=1 and @dist > 3500 and @dist <= 4000, CsDamage.dmg_health, 0)) as d4000_body_health,
    sum(if(hitgroup=1 and @dist > 3500 and @dist <= 4000, CsDamage.dmg_armor, 0)) as d4000_head_armor,
    sum(if(not hitgroup=1 and @dist > 3500 and @dist <= 4000, CsDamage.dmg_armor, 0)) as d4000_body_armor,
    -- 4500
    sum(if(hitgroup=1 and @dist > 4000 and @dist <= 4500, CsDamage.dmg_health, 0)) as d4500_head_health,
    sum(if(not hitgroup=1 and @dist > 4000 and @dist <= 4500, CsDamage.dmg_health, 0)) as d4500_body_health,
    sum(if(hitgroup=1 and @dist > 4000 and @dist <= 4500, CsDamage.dmg_armor, 0)) as d4500_head_armor,
    sum(if(not hitgroup=1 and @dist > 4000 and @dist <= 4500, CsDamage.dmg_armor, 0)) as d4500_body_armor,
    -- 5000
    sum(if(hitgroup=1 and @dist > 4500 and @dist <= 5000, CsDamage.dmg_health, 0)) as d5000_head_health,
    sum(if(not hitgroup=1 and @dist > 4500 and @dist <= 5000, CsDamage.dmg_health, 0)) as d5000_body_health,
    sum(if(hitgroup=1 and @dist > 4500 and @dist <= 5000, CsDamage.dmg_armor, 0)) as d5000_head_armor,
    sum(if(not hitgroup=1 and @dist > 4500 and @dist <= 5000, CsDamage.dmg_armor, 0)) as d5000_body_armor,
    -- inf
    sum(if(hitgroup=1 and @dist > 5000, CsDamage.dmg_health, 0)) as dinf_head_health,
    sum(if(not hitgroup=1 and @dist > 5000, CsDamage.dmg_health, 0)) as dinf_body_health,
    sum(if(hitgroup=1 and @dist > 5000, CsDamage.dmg_armor, 0)) as dinf_head_armor,
    sum(if(not hitgroup=1 and @dist > 5000, CsDamage.dmg_armor, 0)) as dinf_body_armor

  from
    CsDamage

  join
    CsRound
    on CsRound.replay_id in (?)
    and CsRound.id = CsDamage.round_id

  where
    CsDamage.given_eid in (?)

  and
    -- do not include self-inflicted damage
    CsDamage.taker_eid != CsDamage.giver_eid

  group by
    CsRound.replay_id, CsDamage.given_eid

  having
    -- neither position must be null
    distance is not NULL;
  `

	query, args, err := sqlx.In(stmt, matches, players)
	if err != nil {
		return features, fmt.Errorf("Unable to create query: %w", err)
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return features, fmt.Errorf("Unable to execute query: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var f svm.PlayerFeatures
		var match_eid, player_eid structs.Eid
		var temp sql.NullInt64 // throw away distance variable used for calculations
		err = rows.Scan(
			&match_eid,
			&player_eid,

			&temp,

			&f.DamageGivenHealth500Head,
			&f.DamageGivenHealth500Body,
			&f.DamageGivenArmor500Head,
			&f.DamageGivenArmor500Body,

			&f.DamageGivenHealth1000Head,
			&f.DamageGivenHealth1000Body,
			&f.DamageGivenArmor1000Head,
			&f.DamageGivenArmor1000Body,

			&f.DamageGivenHealth1500Head,
			&f.DamageGivenHealth1500Body,
			&f.DamageGivenArmor1500Head,
			&f.DamageGivenArmor1500Body,

			&f.DamageGivenHealth2000Head,
			&f.DamageGivenHealth2000Body,
			&f.DamageGivenArmor2000Head,
			&f.DamageGivenArmor2000Body,

			&f.DamageGivenHealth2500Head,
			&f.DamageGivenHealth2500Body,
			&f.DamageGivenArmor2500Head,
			&f.DamageGivenArmor2500Body,

			&f.DamageGivenHealth3000Head,
			&f.DamageGivenHealth3000Body,
			&f.DamageGivenArmor3000Head,
			&f.DamageGivenArmor3000Body,

			&f.DamageGivenHealth3500Head,
			&f.DamageGivenHealth3500Body,
			&f.DamageGivenArmor3500Head,
			&f.DamageGivenArmor3500Body,

			&f.DamageGivenHealth4000Head,
			&f.DamageGivenHealth4000Body,
			&f.DamageGivenArmor4000Head,
			&f.DamageGivenArmor4000Body,

			&f.DamageGivenHealth4500Head,
			&f.DamageGivenHealth4500Body,
			&f.DamageGivenArmor4500Head,
			&f.DamageGivenArmor4500Body,

			&f.DamageGivenHealth5000Head,
			&f.DamageGivenHealth5000Body,
			&f.DamageGivenArmor5000Head,
			&f.DamageGivenArmor5000Body,

			&f.DamageGivenHealthInfHead,
			&f.DamageGivenHealthInfBody,
			&f.DamageGivenArmorInfHead,
			&f.DamageGivenArmorInfBody,
		)
		if err != nil {
			return features, fmt.Errorf("Unable to scan row of result set: %w", err)
		}

		if features[match_eid] == nil {
			features[match_eid] = make(map[structs.Eid]svm.PlayerFeatures)
		}
		features[match_eid][player_eid] = f

	}
	if rows.Err() != nil {
		return features, fmt.Errorf("An error occurred when processing the result set: %w", err)
	}

	return features, nil
}
