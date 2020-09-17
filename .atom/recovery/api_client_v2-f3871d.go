package queries

import (
	"database/sql"
	"fmt"

	"gitlab.com/abios/user-svc/database"
	"gitlab.com/abios/user-svc/structs"
)

// GetClientV2 fetches the v2 API client of the specified customer
func GetClientV2(customerId int64) (*structs.ClientV2, error) {
	var id, secret string
	var prod bool
	var createdAt []uint8
	var hForw, hBack *int64
	var sReq, mReq, search int64
	// Prepare and execute SQL statement
	stmt := `
		SELECT
			id,
			secret,
			is_prod,
			created_at,
			hours_forward,
			hours_back,
			second_rate,
			minute_rate,
			search_minute_rate
		FROM
			oauth.Client
		WHERE customer_id = ?`
	err := database.OAuthDB.QueryRow(stmt, customerId).Scan(&id, &secret, &prod, &createdAt, &hForw, &hBack, &sReq, &mReq, &search)
	if err != nil {
		return nil, err
	}
	t, err := database.Timestamp2unix(createdAt)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse created_at field: %w", err)
	}
	return &structs.ClientV2{
		Id:         id,
		Secret:     secret,
		IsProd:     prod,
		CustomerId: customerId,
		CreatedAt:  *t,
		Config: structs.ProductConfig{
			TimeAccess: structs.TimeAccess{
				HoursForward: hForw,
				HoursBack:    hBack,
			},
			RequestRates: structs.ReqRates{
				Second: sReq,
				Minute: mReq,
				Search: search,
			},
		},
	}, nil
}

// GetSubscriptionV2 creates a subscription summary for the v2 API client of the
// specified customer
func GetSubscriptionV2(customerId int64) (*structs.SubscriptionV2, error) {
	rows, err := database.OAuthDB.Query(`
		SELECT distinct
		  game_id,
		  level,
		  X.client_id
		FROM
		  (
		    select
		      client_id
		    FROM
		      dashboard.users_2_client_id
		    join
		      dashboard.Users
		    on
		      dashboard.users_2_client_id.email = dashboard.Users.email
		    join
		      dashboard.User2Customer
		    on
		      dashboard.Users.id = dashboard.User2Customer.user_id
		    where
		      dashboard.User2Customer.customer_id = 1
		    limit 1
		  ) as X
		join
		  oauth.ClientScope on oauth.ClientScope.client_id = X.client_id
		join
		  oauth.GameScope on oauth.GameScope.client_id = X.client_id;
	`)
	if err != nil {
		return []structs.Product{}, err
	}

	var products []structs.Product = []structs.Product{}
	var id sql.NullInt64
	var stripeId sql.NullString

	// Process each row of the result-set
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &stripeId); err != nil {
			return products, err
		}
		products = append(products, structs.Product{
			Id:       id.Int64,
			StripeId: stripeId.String,
		})
	}
	return products, nil
}
