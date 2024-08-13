package mysql

import (
	"fmt"

	"github.com/armiariyan/assessment-tsel/internal/config"

	"github.com/armiariyan/bepkg/database/mysql"
)

func NewMySQL(db config.DB) (client *mysql.SQLDB) {
	opts := &mysql.Options{
		DSN:                db.URI,
		MinIdleConnections: db.MinPoolSize,
		MaxOpenConnections: db.MaxPoolSize,
		MaxLifetime:        db.Timeout,
		LogMode:            db.DebugMode,
	}

	client, err := mysql.Connect("1", opts)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to mysql: %v", err))
	}

	fmt.Println("[SUCCESS] connect to DB MYSQL .........")

	return
}
