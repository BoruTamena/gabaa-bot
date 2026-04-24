package initiator

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

func InitViper(currentDir string) error {

	viper.AddConfigPath(currentDir + "/config")
	// viper.AddConfigPath("config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Print("config file not found, using environment variables")
		} else {
			log.Print("failed to read config", err)
			return err
		}
	}

	return nil

}
