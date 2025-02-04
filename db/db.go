package db

import (
	"fmt"
	"movie_management/config"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql" 
)


func Connect() error {
	cfg, err := cnf.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		cfg.Mysql.DB_USER, cfg.Mysql.DB_PASSWORD, cfg.Mysql.DB_HOST, cfg.Mysql.DB_PORT, cfg.Mysql.DB_NAME)

	if err := orm.RegisterDataBase("default", "mysql", connectionString); err != nil {
		return fmt.Errorf("failed to register database: %v", err)
	}

	orm.Debug = true 
	return orm.RunSyncdb("default", false, true) 
}

// func GetDB() orm.Ormer {
// 	return orm.NewOrm()
// }