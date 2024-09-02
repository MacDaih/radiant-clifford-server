package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"webservice/config"
	"webservice/internal/handler"
	"webservice/internal/repository"
	"webservice/internal/service"

	"webservice/cmd/server"
	"webservice/cmd/worker"

	"github.com/gorilla/mux"
	sdk "github.com/macdaih/porter_go_sdk"
)

func main() {

	appCTX, shutdown := context.WithCancel(context.Background())
	config.Boot()

	log.Println("starting webservice")

	repo := repository.NewReportRepository(config.GetDBEnv())

	recordReport := service.RecordReportFunc(repo)
	hdlr := handler.NewServiceHandler(repo)

	httpError := make(chan error, 1)
	sysInt := make(chan os.Signal, 1)

	router := mux.NewRouter()

	router.HandleFunc("/reports/{range}", hdlr.GetReportsFrom).Methods("GET")
	router.HandleFunc("/by_date/{date}", hdlr.GetReportsByDate).Methods("GET")

	webservice := http.Server{
		Addr:    config.GetPort(),
		Handler: router,
	}

	go server.RunWebservice(&webservice, httpError)

	client := sdk.NewClient(
		config.GetServerAddr(),
		15,
		sdk.WithID(config.GetClientID()),
		sdk.WithBasicCredentials(config.GetUserName(), config.GetUserPasswd()),
		sdk.WithCallBack(
			recordReport,
		),
	)

	consumer := worker.NewConsumer(client, config.GetTopics())

	go func() {
		consumer.Run(appCTX)
	}()

	signal.Notify(sysInt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-httpError:
		shutdown()
		if err != nil {
			log.Fatalf("http Server error : %s", err.Error())
		}
	case <-sysInt:
		shutdown()
		// TODO move webservice shutdown into its own package
		log.Println("interrupt : service is shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		defer cancel()

		if err := webservice.Shutdown(ctx); err != nil {
			log.Printf("error when shutting down server : %s", err.Error())
			log.Println("closing webservice ...")
			webservice.Close()
		}

		return
	}
}
