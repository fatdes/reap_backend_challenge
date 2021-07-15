package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fatdes/reap_backend_challenge/log"
	"github.com/fatdes/reap_backend_challenge/route"
)

func main() {
	ctx := context.Background()
	log := log.GetLogger("main")
	t := time.Now()
	defer func() {
		log.Add("canonical", "yes").Add("uptime_seconds", time.Now().Sub(t).Seconds()).Info("main")
		log.Sync()
	}()
	router := route.NewRouter()

	port := 8080
	if portString, found := os.LookupEnv("PORT"); found {
		newPort, err := strconv.Atoi(portString)
		if err != nil {
			log.Fatal(fmt.Sprintf("invalid PORT=%s", portString), err)
		}
		port = newPort
	}
	ws := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		log.Info("starting web service")
		if err := ws.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("fail to start web server", err)
		}
	}()

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info("stopping everything by signal ", <-termChan)

	log.Info("stopping web server")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := ws.Shutdown(ctx); err != nil {
		log.Error("fail to stop web server", err)
	}
	log.Info("stopped web server")

	log.Info("stopped everything")
}
