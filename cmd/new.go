package cmd

import (
	"github.com/fatih/color"
	"github.com/reo7sp/two-step-migrate/migrationcreator"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"regexp"
	"strings"
)

var newCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n"},
	Short:   "Create a new migration",
	Run: func(cmd *cobra.Command, args []string) {
		questions := []*survey.Question{
			{
				Name:     "name",
				Prompt:   &survey.Input{Message: "Enter name of the new migration:"},
				Validate: survey.Required,
				Transform: survey.TransformString(func(s string) string {
					return regexp.MustCompile("[^a-z]").ReplaceAllString(strings.ToLower(s), "_")
				}),
			},
		}
		answers := struct {
			Name string
		}{}

		err := survey.Ask(questions, &answers)
		if err != nil {
			log.Fatal(err)
		}

		err = migrationcreator.Create(answers.Name)
		if err != nil {
			log.Fatal(err)
		}
		color.Green("Created new migration %s", answers.Name)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
