package main

import (
	"fmt"
	"net/http"
	"github.com/best-project/api/internal/server"
	"github.com/sirupsen/logrus"
	"github.com/best-project/api/internal/storage"
	"github.com/best-project/api/internal/config"
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
