package main

import (
	"fmt"
	"net/http"
	"database/sql"
	"github.com/best-project/api/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("hello world")

	_, err := sql.Open("mysql", "root:1234@/litdb")
	fatalOnError(err)

	srv := server.NewServer()
	fatalOnError(http.ListenAndServe(":8080", srv.Handle()))
}

func fatalOnError(err error) {
	if err != nil {
		// change to Fatal
		logrus.Print(err.Error())
	}
}
