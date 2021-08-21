package main

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/wilhelmeek/dayforit"
)

func main() {
	logger := logrus.New()

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/check", dayforit.Check)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.WithError(err).Info("handling http traffic")
	}
}
