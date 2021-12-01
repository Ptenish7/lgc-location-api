package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ozonmp/lgc-location-api/internal/retranslator/retranslator"
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
	defer rt.Close()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
