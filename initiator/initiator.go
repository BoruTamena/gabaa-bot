package initiator

import (
	"fmt"

	"github.com/BoruTamena/gabaa-bot/internal/constant/persistencedb"
	"github.com/spf13/viper"
)

func Init() {

	err := InitViper("./")

	if err != nil {
		panic(err)
	}

	db_peristsence := persistencedb.NewPersistenceDb()

	// creating migiration
	mg := InitMigiration(viper.GetString("migiration.path"),
		viper.GetString("db.url"))

	persistence := InitPersistence(db_peristsence)
	UpMigiration(mg)

	fmt.Println("migration created")
	platform := InitPlatFormLayer()

	group := platform.tg.Group()

	module := InitModule(persistence, platform)

	handler := InitHandler(module)

	InitRoute(&group, handler)

	// starting bot
	platform.tg.Start()

	fmt.Println("bot listing ...")

}
