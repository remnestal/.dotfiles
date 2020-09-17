package queries

import (
	"database/sql"
	"encoding/json"

	"github.com/sirupsen/logrus"

	"gitlab.com/abios/user-svc/database"
	"gitlab.com/abios/user-svc/structs"
)

// GetDefaultProduct fetches information about the specified product from the
// database. The JSON content of the config-attribute is unmarshaled as well,
// meaning that a JSON parsing error may occur if that data is not formatted
// properly. The DefaultPlan field is not set by this function and must be
// separately fetched via the Stripe SDK
func GetDefaultProduct(pid int64, logger *logrus.Entry) (structs.DefaultProduct, error) {
	var id sql.NullInt64
	var name, config sql.NullString
	var stripeId *string
	// Prepare and execute SQL statement
	stmt := `
		SELECT
			DefaultProductConfig.id,
			DefaultProductConfig.name,
			DefaultProductConfig.config,
			StripeProductMapping.stripe_product_id
		FROM
			DefaultProductConfig
		LEFT JOIN StripeProductMapping ON
			StripeProductMapping.id = DefaultProductConfig.stripe_product
		WHERE
			DefaultProductConfig.id = ?
	`
	err := database.OAuthDB.QueryRow(stmt, pid).Scan(&id, &name, &config, &stripeId)
	if err != nil {
		logger.WithError(err).WithField("sql-statement", stmt).Errorln("Unable to query database for products")
		return structs.Product{}, err
	}
	// Unmarshal the JSON config field
	var productConfig structs.ProductConfig
	if err := json.Unmarshal([]byte(config.String), &productConfig); err != nil {
		logger.WithError(err).WithField("json", config.String).Errorln("Unable to unmarshal product configuration")
		return structs.Product{}, err
	}
	return structs.DefaultProduct{
		Id:       id.Int64,
		Name:     name.String,
		Config:   productConfig,
		StripeId: stripeId,
	}, nil
}

func GetAllProducts() ([]structs.Product, error) {
	// Prepare and execute SQL statement
	stmt := `
		SELECT
			id,
			stripe_product_id
		FROM
			StripeProductMapping
	`
	rows, err := database.DashbordDB.Query(stmt, uid)
	if err != nil {
		return []structs.Product, err
	}
	defer rows.Close()
	// Process each row of the result-set
	var products []structs.Product = []structs.Product{}
	var id sql.NullInt64
	var stripeId sql.NullString
	for rows.Next() {
		if err := rows.Scan(&id, &stripeId); err != nil {
			return roles, err
		}
		products = append(products, structs.Product{
			CustomerId: customer_id.Int64,
			RoleId:     role_id.Int64,
		})
	}
	return products, nil
}
