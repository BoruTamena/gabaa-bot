package initiator

import (
	"log"
	"strings"

	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/spf13/viper"
)

func InitViper(currentDir string) error {

	// Read from .env file in the project root
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(currentDir)

	// Map dot-separated viper keys (e.g. "db.url") to
	// underscore-separated env vars (e.g. DB_URL)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Print(".env file not found, falling back to environment variables only")
		} else {
			log.Printf("failed to read .env file: %v", err)
			return err
		}
	}

	logger.Info("config loaded successfully")

	return nil

}
