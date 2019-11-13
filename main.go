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

	fb := facebook.New()
	fb.SetAppID(cfg.FbAppKey)
	fb.SetAppSecret(cfg.FbAppSecret)

	db, err := storage.NewDatabase(cfg)
	fatalOnError(err)

	srv := server.NewServer(db, fb)
	fatalOnError(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), srv.Handle()))
}

func fatalOnError(err error) {
	if err != nil {
		logrus.Fatal(err.Error())
	}
}
