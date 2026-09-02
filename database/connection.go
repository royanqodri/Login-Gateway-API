package database

import (
	"fmt"
	"time"

	"github.com/royanqodri/Login-Gateway-API/config"
	"github.com/royanqodri/Login-Gateway-API/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConnMaster    *gorm.DB
	DBConnOperation *gorm.DB
)

// ConnectDatabase establishes a connection to the master database.
func ConnectDatabase() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Get().DBLoginConfig.Host,
		config.Get().DBLoginConfig.User,
		config.Get().DBLoginConfig.Password,
		config.Get().DBLoginConfig.DatabaseName,
		config.Get().DBLoginConfig.Port,
		config.Get().DBLoginConfig.SslMode,
		config.Get().DBLoginConfig.Timezone,
	)
	DBConnMaster, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	util.PanicIfError(err)

	// Get generic database object sql.DB to use its functions
	sqlDB, err := DBConnMaster.DB()
	util.PanicIfError(err)

	sqlDB.SetMaxIdleConns(config.Get().DBLoginConfig.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.Get().DBLoginConfig.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(config.Get().DBLoginConfig.MaxConnLifetime * time.Minute)
	sqlDB.SetConnMaxIdleTime(config.Get().DBLoginConfig.MaxConnIdletime * time.Minute)

	// Force apply statement_timeout to every new connection in pool
	sqlDB.Exec("SET statement_timeout = '15s'") // <- will apply only to this session
}
