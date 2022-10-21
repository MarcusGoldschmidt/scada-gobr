package main

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	err := pkg.ConfigureFlags()
	if err != nil {
		log.Panic(err)
	}

	opt, err := pkg.ParseOptions()

	if err != nil {
		log.Panic(err)
	}

	scada, err := pkg.DefaultScadagobr(opt)

	if err != nil {
		log.Panic(err)
	}
	ctx := context.Background()

	err = scada.Setup(ctx)

	if err != nil {
		log.Panic(err)
	}

	err = scada.Run(ctx)
	if err != nil {
		log.Panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(ctx, opt.ShutdownWait)
	defer cancel()

	scada.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
