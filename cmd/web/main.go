package main

import (
	"crypto/tls"
	"forum/app"
	"forum/internal/config"
	"forum/internal/handlers"
	"forum/internal/repo"
	"forum/internal/service"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conf := config.Loader()

	templateCache, err := app.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := app.New(infoLog, errorLog, templateCache)

	repo, err := repo.New(conf.StoragePath)
	if err != nil {
		errorLog.Fatal(err)
	}

	serv := service.New(repo)

	hand := handlers.New(serv, app, conf)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS12,
	}

	tlsCert := "./tls/cert.pem"
	tlsKey := "./tls/key.pem"

	srv := &http.Server{
		Addr:         conf.Address,
		ErrorLog:     errorLog,
		Handler:      hand.Routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on https://localhost%s/", conf.Address)
	err = srv.ListenAndServeTLS(tlsCert, tlsKey)
	errorLog.Fatal(err)
}
