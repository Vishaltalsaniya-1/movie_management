package db

import (
	"fmt"
	"movie_management/config"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql" 
)

func Connect() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME)

	err = orm.RegisterDataBase("default", "mysql", connectionString)
	if err != nil {
		return fmt.Errorf("failed to register database: %v", err)
	}

	orm.Debug = true

	return orm.RunSyncdb("default", false, true)
}

func GetDB() orm.Ormer {
	return orm.NewOrm()
}
