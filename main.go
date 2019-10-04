package main

import (
	"fmt"
	"github.com/best-project/api/internal/config"
	"github.com/best-project/api/internal/server"
	"github.com/best-project/api/internal/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	cfg := config.NewConfig()

	db, err := storage.NewDatabase(cfg)
	fatalOnError(err)

	srv := server.NewServer(db)
	fatalOnError(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), srv.Handle()))
}

func fatalOnError(err error) {
	if err != nil {
		// TODO: change to Fatal
		logrus.Print(err.Error())
	}
}
