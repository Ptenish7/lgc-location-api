package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ozonmp/lgc-location-api/internal/app/retranslator"
)

func main() {
	sigs := make(chan os.Signal, 1)

	cfg := retranslator.Config{
		ChannelSize:   512,
		ConsumerCount: 2,
		ConsumerSize:  10,
		ProducerCount: 28,
		WorkerCount:   2,
		BatchSize:     1,
	}

	rt := retranslator.NewRetranslator(cfg)
	rt.Start()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
