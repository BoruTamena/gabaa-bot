package persistencedb

import (
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PersistenceDb struct {
	*gorm.DB
}

func NewPersistenceDb() PersistenceDb {

	url := viper.GetString("db.url")

	g_db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {

		logger.Error("Failed to connect to database",
			zap.String("error", err.Error()),
			zap.String("url", url))
		panic(err)
	}
	return PersistenceDb{
		g_db,
	}
}
