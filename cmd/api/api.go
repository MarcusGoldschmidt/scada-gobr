package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"scadagobr/pkg"
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

	err = scada.Run(ctx)
	if err != nil {
		log.Panic(err)
	}
}
