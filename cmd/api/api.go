package main

import (
	"context"
	api "github.com/MarcusGoldschmidt/scadagobr/cmd/api/command"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	ctx := context.Background()

	err := api.Execute(ctx)
	if err != nil {
		log.Panic(err)
	}
	os.Exit(0)
}
