package postgresql

import (
	"fmt"
	"log"

	"github.com/armiariyan/assessment-tsel/internal/config"
	"github.com/labstack/gommon/color"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg config.PostgresqlDB) (db *gorm.DB) {
	var dsn string
	if config.GetString("env") == "local" {
		// * if somehow use local db without password, remove the string "password=%s" below
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)
	} else {
		// * for assessment purpose, use psql in railway directly
		dsn = "postgresql://postgres:oZULeLKKzjLBskxHyfdTkzmRDkKtUbzz@roundhouse.proxy.rlwy.net:54482/railway"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if cfg.Debug {
		db = db.Debug()
	}

	color.Println(color.Green(fmt.Sprintf("⇨ connected to postgresql db on %s\n", cfg.Name)))
	return
}
