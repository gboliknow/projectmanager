package main

import (
	"log"
	"projectmanager/api"
	"projectmanager/internal/config"
	"projectmanager/internal/database"

	"github.com/go-sql-driver/mysql"
)

func main() {

	cfg := mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	sqlStorage := database.NewMySQLStorage(cfg)

	db, err := sqlStorage.InitializeDatabase()

	if err != nil {
		log.Fatal(err)
	}
	store := api.NewStore(db)
	api := api.NewAPIServer(":3306", store)
	api.Serve()
}
