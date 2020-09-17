package queries

import (
	"database/sql"

	"gitlab.com/abios/user-svc/database"
	"gitlab.com/abios/user-svc/structs"
)

func GetCardByFingerprint(fingerprint string) (*CardMapping, error) {
  var customerId sql.NullInt64
	var fingerprint sql.NullString
	// Prepare and execute SQL statement
	stmt := `
    SELECT
      customer_id,
      stripe_fingerprint
    FROM
      oauth.StripeCardFingerprint
    WHERE
      stripe_fingerprint = ?
    AND
      deleted_at IS NULL`
	err := database.OAuthDB.QueryRow(stmt, fingerprint).Scan(&id, &customerId, &fingerprint)
  if err == sql.ErrNoRows {
    return nil, nil
  }
	if err != nil {
		return nil, err
	}
	return structs.Customer{
		Id:       id.Int64,
		Name:     name.String,
		StripeId: stripeId,
	}, nil


}
