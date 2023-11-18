package main

import (
	"context"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FlashpointProject/CommunityWebsite/config"
	"github.com/FlashpointProject/CommunityWebsite/database"
	"github.com/FlashpointProject/CommunityWebsite/logging"
	"github.com/FlashpointProject/CommunityWebsite/service"
	"github.com/FlashpointProject/CommunityWebsite/transport"
	"github.com/FlashpointProject/CommunityWebsite/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("IN_KUBERNETES") != "" {
		_, err := os.Stat(".env")
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}

	conf, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}

	log := logging.InitLogger()
	l := log.WithField("fpcomm", map[string]interface{}{"version": conf.Version})
	l.Infoln("Starting server...")

	router := mux.NewRouter()
	mime.AddExtensionType(".js", "application/javascript")

	pgdb := database.OpenPostgresDB(l, conf)
	defer pgdb.Close()

	fpfss, err := transport.NewFpfss(conf.OauthConfig, conf.FpfssApiUrl)
	if err != nil {
		l.WithError(err).Fatalln("failed to connect to fpfss")
	}
	app := &transport.App{
		Conf:    conf,
		Service: service.NewService(pgdb, conf.SessionExpirationSeconds),
		CC: utils.CookieCutter{
			Previous: securecookie.New([]byte(conf.SecurecookieHashKeyPrevious), []byte(conf.SecurecookieBlockKeyPrevious)),
			Current:  securecookie.New([]byte(conf.SecurecookieHashKeyCurrent), []byte(conf.SecurecookieBlockKeyPrevious)),
		},
		Fpfss: fpfss,
	}

	err = app.Service.LoadRoles(context.Background())
	if err != nil {
		l.WithError(err).Fatalln("failed to load roles")
	}

	srv := &http.Server{
		Handler:      logging.LogRequestHandler(l, app.Fpfss.WithFpfss(router)),
		Addr:         fmt.Sprintf("0.0.0.0:%d", conf.Port),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	go func() {
		app.ServeRouter(l, srv, router)
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-term
	l.Infoln("signal received, exitting")

	l.Infoln("shutting down the server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		l.WithError(err).Errorln("server shutdown failed")
	}

}
