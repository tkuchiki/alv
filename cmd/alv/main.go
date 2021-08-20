package main

import (
	"log"

	"github.com/tkuchiki/alv/cmd/alv/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
