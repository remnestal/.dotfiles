package database

import (
  "os"
  "fmt"
  "database/sql"

  _ "github.com/go-sql-driver/mysql"
)

const DSN_VAR string = "FORECAST_DSN"

var ForecastDB *sql.DB

func Setup() error {
  if dsn := os.Getenv(DSN_VAR); dsn == "" {
    return fmt.Errorf("DSN environment variable (%v) was not set.", DSN_VAR)
  }
}
