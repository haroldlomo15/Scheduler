package main

import (
	"github.com/haroldlomo15/scheduler/pkg/db"
	"github.com/haroldlomo15/scheduler/pkg/router"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	ok := run()
	if !ok {
		os.Exit(1)
	}
}

func run() bool {
	err := db.New()
	if err != nil {
		log.Error(err)
		return false
	}

	log.Info("Listening on port 7005...")
	err = http.ListenAndServe("0.0.0.0:7005", router.NewRouter())
	if err != nil {
		log.Errorf("error listenAndServe - %v", err)
		return false
	}
	return true
}
