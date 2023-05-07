package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asiman161/re-platform/app/bootstrap"

	"github.com/caarlos0/env/v6"
)

func main() {
	cfg := bootstrap.AppConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := bootstrap.SetupServer(app, cfg)
	startServer(server)
}

func startServer(server *http.Server) {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func(sig chan os.Signal) {
		<-sig
		log.Println("received shutdown event")

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out... forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}(sig)

	log.Printf("start server at: %v\n", server.Addr)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}
