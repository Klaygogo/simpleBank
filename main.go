package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"

	api "github.com/Klaygogo/simplebank/api"
	db "github.com/Klaygogo/simplebank/db/sqlc"
	"github.com/Klaygogo/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("无法加载配置:", err)
	}
	pool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("无法创建数据库连接池:", err)
	}
	store := db.NewStore(pool)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("无法启动服务器:", err)
	}

}
