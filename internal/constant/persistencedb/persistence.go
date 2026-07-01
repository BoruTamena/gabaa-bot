package persistencedb

import (
	"github.com/spf13/viper"
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
		panic(err)
	}
	return PersistenceDb{
		g_db,
	}
}
