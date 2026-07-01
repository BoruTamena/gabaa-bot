package initiator

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func InitViper(currentDir string) error {

	// Load .env into the process environment.
	// godotenv.Load is a no-op (logs a warning) if the file is missing,
	// so the app still works when env vars are injected by Docker / k8s / CI.
	envFile := filepath.Join(currentDir, ".env")
	if err := godotenv.Load(envFile); err != nil {
		log.Printf(".env file not found at %s, relying on OS environment variables: %v", envFile, err)
	}

	// Map dot-separated viper keys (e.g. "db.url") → underscore env vars (e.g. "DB_URL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	logger.Info("config loaded successfully")

	return nil
}
