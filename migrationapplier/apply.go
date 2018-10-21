package migrationapplier

import (
	"database/sql"
	"github.com/pkg/errors"
	"io/ioutil"
)

func ApplyUpMigration(m Migration, db *sql.DB) error {
	return applyMigration(m, true, db)
}

func ApplyDownMigration(m Migration, db *sql.DB) error {
	return applyMigration(m, false, db)
}

func applyMigration(migration Migration, isApply bool, db *sql.DB) error {
	var filePath string
	if isApply {
		filePath = migration.UpMigrationPath()
	} else {
		filePath = migration.DownMigrationPath()
	}

	migrationContentsBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.Wrap(err, "cannot load file "+filePath)
	}
	migrationContents := string(migrationContentsBytes)

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "cannot start transaction")
	}
	defer tx.Rollback()

	_, err = tx.Exec(migrationContents)
	if err != nil {
		return errors.Wrap(err, "cannot exec migration sql code")
	}

	if isApply {
		_, err = tx.Exec(`INSERT INTO `+sqlTableName+` (name, kind) VALUES ($1, $2)`, migration.Name, migration.Kind)
	} else {
		_, err = tx.Exec(`DELETE FROM `+sqlTableName+` WHERE name = $1 AND kind = $2`, migration.Name, migration.Kind)
	}
	if err != nil {
		return errors.Wrap(err, "cannot save state")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "cannot commit transaction")
	}

	return nil
}
