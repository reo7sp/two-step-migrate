package migrationcreator

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/reo7sp/two-step-migrate/migrationcommon"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func Create(name string) error {
	migrationName := generateName(name)

	newFolderPath := filepath.Join(migrationcommon.MigrationsFolderName, migrationName)
	err := os.MkdirAll(newFolderPath, 0755)
	if err != nil {
		return errors.Wrap(err, "cannot create migration path")
	}

	filesToCreate := []string{
		migrationcommon.FirstDownMigrationFileName,
		migrationcommon.FirstUpMigrationFileName,
		migrationcommon.SecondDownMigrationFileName,
		migrationcommon.SecondUpMigrationFileName,
	}
	migrationContents := []byte("-- TODO\n")
	for _, it := range filesToCreate {
		newFilePath := filepath.Join(newFolderPath, it)
		err := ioutil.WriteFile(newFilePath, migrationContents, 0644)
		if err != nil {
			return errors.Wrap(err, "cannot create migration file")
		}
	}

	return nil
}

func generateName(name string) string {
	timeStr := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s", timeStr, name)
}
