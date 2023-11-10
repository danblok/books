package main

import (
	"log"

	"github.com/danblok/books/internal/pkg/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalln(err)
	}
	app.Run()
}
