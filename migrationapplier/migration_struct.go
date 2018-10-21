package migrationapplier

import (
	"fmt"
	"github.com/reo7sp/two-step-migrate/migrationcommon"
	"path/filepath"
)

type Migration struct {
	Name string
	Kind string
}

func (m Migration) UpMigrationPath() string {
	var fileName string
	switch m.Kind {
	case migrationcommon.FirstMigrationKind:
		fileName = migrationcommon.FirstUpMigrationFileName
	case migrationcommon.SecondMigrationKind:
		fileName = migrationcommon.SecondUpMigrationFileName
	default:
		panic(fmt.Sprintf("unknown migration kind: %s", m.Kind))
	}
	return filepath.Join(migrationcommon.MigrationsFolderName, m.Name, fileName)
}

func (m Migration) DownMigrationPath() string {
	var fileName string
	switch m.Kind {
	case migrationcommon.FirstMigrationKind:
		fileName = migrationcommon.FirstDownMigrationFileName
	case migrationcommon.SecondMigrationKind:
		fileName = migrationcommon.SecondDownMigrationFileName
	default:
		panic(fmt.Sprintf("unknown migration kind: %s", m.Kind))
	}
	return filepath.Join(migrationcommon.MigrationsFolderName, m.Name, fileName)
}

func (m Migration) IsLess(other Migration) bool {
	if m.Name != other.Name {
		return m.Name < other.Name
	} else {
		return m.Kind < other.Kind
	}
}

type MigrationsSorted []Migration

func (a MigrationsSorted) Len() int {
	return len(a)
}

func (a MigrationsSorted) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a MigrationsSorted) Less(i, j int) bool {
	return a[i].IsLess(a[j])
}
