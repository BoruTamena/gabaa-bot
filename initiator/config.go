package initiator

import (
	"log"

	"github.com/spf13/viper"
)

func InitViper(currentDir string) error {

	viper.AddConfigPath(currentDir + "/config")
	// viper.AddConfigPath("config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Print("failed to read config", err)
		return err
	}

	return nil

}
