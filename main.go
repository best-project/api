package main

import (
	"fmt"
	"github.com/best-project/api/internal/config"
	"github.com/best-project/api/internal/server"
	"github.com/best-project/api/internal/storage"
	"github.com/madebyais/facebook-go-sdk"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()
	fatalOnError(err)

	logger := logrus.New()
	logger.Info("Starting application")

	fb := facebook.New()
	fb.SetAppID(cfg.FbAppKey)
	fb.SetAppSecret(cfg.FbAppSecret)

	db, err := storage.NewDatabase(cfg, logger)
	fatalOnError(err)

	srv := server.NewServer(db, fb, logger)
	logger.Info("===Starting Server===")
	fatalOnError(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), srv.Handle()))
}

func fatalOnError(err error) {
	if err != nil {
		logrus.Fatal(err.Error())
	}
}
