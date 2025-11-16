package repositories

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

var gormDB *gorm.DB

func SetupGormDB(config GormDBConfig) {
	dsn := "host=" + config.Host +
		" user=" + config.User +
		" password=" + config.Password +
		" dbname=" + config.Dbname +
		" port=" + config.Port +
		" sslmode=disable" +
		" TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	gormDB = db
}
