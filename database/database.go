package database

import (
	"encore.dev/storage/sqldb"
)

// Database is the database object for the application.
var Database = sqldb.NewDatabase("database", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
