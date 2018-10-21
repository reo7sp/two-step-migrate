package cfg

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func ParseDatabaseUrl() {
	viper.BindEnv("DATABASE_URL")
	if !viper.IsSet("database_url") {
		fmt.Printf("DATABASE_URL environment variable is required\n")
		os.Exit(2)
	}
	if !strings.Contains(viper.GetString("database_url"), "ssl-mode=") {
		viper.Set("database_url", viper.GetString("database_url")+"?sslmode=disable")
	}
}
