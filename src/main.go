package main

import (
	"fmt"
	"os"
	"os/signal"

	"citadel_intranet/src/application"
	"citadel_intranet/src/config"
	"citadel_intranet/src/db"
	"citadel_intranet/src/server"
)

func main() {
	cfg := config.LoadConfig()
	dbClient := db.NewDatabaseClient(cfg)

	db.Migrate(dbClient.Db, cfg.MigrationsPath)

	webServer := server.NewServer(cfg)
	app := application.NewApp(dbClient, webServer)
	defer app.Close()
	app.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	fmt.Printf("Received signal, shutting down: %s", s)
	app.Close()
}
