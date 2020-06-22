package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blakelead/nsinjector/internal/config"
	"github.com/blakelead/nsinjector/pkg/controller"
	log "github.com/sirupsen/logrus"
)

func main() {

	clients, err := config.NewClients()
	if err != nil {
		log.Fatal(err)
	}

	informerFactories := config.NewFactories(clients, time.Second*15)

	c := controller.NewController(
		clients,
		informerFactories,
	)

	stop := make(chan struct{})
	defer close(stop)

	informerFactories.Start(stop)

	if err := c.Run(1, stop); err != nil {
		log.Fatalf("Error running controller: %s", err.Error())
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}
