package cmd

import (
	"database/sql"
	"github.com/fatih/color"
	_ "github.com/lib/pq"
	"fmt"
	"github.com/reo7sp/two-step-migrate/migrationapplier"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"strings"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Aliases: []string{"u"},
	Short: "Apply migration",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("postgres", viper.GetString("database_url"))
		if err != nil {
			log.Fatal(err)
		}

		err = migrationapplier.EnsureBootstrap(db)
		if err != nil {
			log.Fatal(err)
		}

		migrations, err := migrationapplier.ListNextUpMigrations(db)
		if err != nil {
			log.Fatal(err)
		}
		if len(migrations) == 0 {
			color.Green("No migrations available to apply")
			return
		}

		migrationNamesForSurvey := make([]string, 0, len(migrations))
		for _, m := range migrations {
			migrationNamesForSurvey = append(migrationNamesForSurvey, fmt.Sprintf("%s %s", m.Name, m.Kind))
		}

		questions := []*survey.Question{
			{
				Name:     "names",
				Prompt:   &survey.MultiSelect{
					Message: "Select migrations to apply:",
					Options: migrationNamesForSurvey,
				},
				Validate: survey.Required,
			},
		}
		answers := struct {
			Names []string
		}{}

		err = survey.Ask(questions, &answers)
		if err != nil {
			log.Fatal(err)
		}

		migrationsToApply := make(migrationapplier.MigrationsSorted, 0, len(answers.Names))
		for _, s := range answers.Names {
			parts := strings.SplitN(s, " ", 2)
			migrationsToApply = append(migrationsToApply, migrationapplier.Migration{Name: parts[0], Kind: parts[1]})
		}

		for _, m := range migrationsToApply {
			log.Printf("Applying %s %s\n", m.Name, m.Kind)
			err := migrationapplier.ApplyUpMigration(m, db)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Done     %s %s\n", m.Name, m.Kind)
		}
		color.Green("Success")
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
