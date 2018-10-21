package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "two-step-migrate",
	Short: "SQL migrate tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.BindEnv("DATABASE_URL")
	if !viper.IsSet("database_url") {
		fmt.Printf("DATABASE_URL environment variable is required\n")
		os.Exit(2)
	}
	if !strings.Contains(viper.GetString("database_url"), "ssl-mode=") {
		viper.Set("database_url", viper.GetString("database_url")+"?sslmode=disable")
	}
}
