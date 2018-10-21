package migrationapplier

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/reo7sp/two-step-migrate/migrationcommon"
	"io/ioutil"
	"sort"
)

func ListAllMigrations() (MigrationsSorted, error) {
	files, err := ioutil.ReadDir(migrationcommon.MigrationsFolderName)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list migrations")
	}

	result := make(MigrationsSorted, 0, len(files)*2)
	for _, f := range files {
		result = append(result, Migration{Name: f.Name(), Kind: migrationcommon.FirstMigrationKind})
		result = append(result, Migration{Name: f.Name(), Kind: migrationcommon.SecondMigrationKind})
	}

	sort.Sort(result)

	return result, nil
}

func ListAppliedMigrations(db *sql.DB) (MigrationsSorted, error) {
	rows, err := db.Query(`SELECT name, kind FROM ` + sqlTableName)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query migrations")
	}
	defer rows.Close()

	result := make(MigrationsSorted, 0)
	for rows.Next() {
		m := Migration{}
		err := rows.Scan(&m.Name, &m.Kind)
		if err != nil {
			return nil, errors.Wrap(err, "failed scan")
		}
		result = append(result, m)
	}

	sort.Sort(result)

	return result, nil
}

func ListNextUpMigrations(db *sql.DB) (MigrationsSorted, error) {
	return listNextMigrations(false, db)
}

func ListNextDownMigrations(db *sql.DB) (MigrationsSorted, error) {
	return listNextMigrations(true, db)
}

func listNextMigrations(mustBeApplied bool, db *sql.DB) (MigrationsSorted, error) {
	allMigrations, err := ListAllMigrations()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get all migrations paths")
	}

	appliedMigrations, err := ListAppliedMigrations(db)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get applied migrations names")
	}

	result := make(MigrationsSorted, 0)
	i := 0
	j := 0
	for i < len(allMigrations) && j < len(appliedMigrations) {
		if allMigrations[i] == appliedMigrations[j] {
			if mustBeApplied {
				result = append(result, allMigrations[i])
			}
			i += 1
			j += 1
		} else {
			if allMigrations[i].IsLess(appliedMigrations[j]) {
				if !mustBeApplied {
					result = append(result, allMigrations[i])
				}
				i += 1
			} else {
				j += 1
			}
		}
	}
	if mustBeApplied {
		// nothing
	} else {
		if i < len(allMigrations) {
			result = append(result, allMigrations[i:]...)
		}
	}

	return result, nil
}
