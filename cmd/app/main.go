package main

import (
	"REST_API_WITH_GO/api"
	"REST_API_WITH_GO/internal/config"
	"REST_API_WITH_GO/internal/database"
	"log"

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

	db, err := sqlStorage.Init()

	if err != nil {
		log.Fatal(err)
	}
	store := api.NewStore(db)
	api := api.NewAPIServer(":3306", store)
	api.Serve()
}
