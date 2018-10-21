package migrationapplier

import (
	"database/sql"
	"github.com/pkg/errors"
)

const sqlTableName = `schema_migrations`

func EnsureBootstrap(db *sql.DB) error {
	var isBootstrapped bool
	checkBootstrapQuery := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables WHERE table_name = $1
		);
	`
	row := db.QueryRow(checkBootstrapQuery, sqlTableName)
	err := row.Scan(&isBootstrapped)
	if err != nil {
		return errors.Wrap(err, "cannot check bootstrap state")
	}

	if !isBootstrapped {
		const doBootstrapQuery = `
			CREATE TABLE schema_migrations (
				name VARCHAR(255),
				kind VARCHAR(16),
				PRIMARY KEY (name, kind)
			);
		`
		_, err := db.Exec(doBootstrapQuery)
		if err != nil {
			return errors.Wrap(err, "cannot create "+sqlTableName+" table")
		}
	}

	return nil
}
