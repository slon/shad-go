// +build !change

package main

import (
	"flag"

	"gitlab.com/slon/shad-go/todo/app"
	"gitlab.com/slon/shad-go/todo/models"
)

func main() {
	port := flag.Int("port", 8080, "port to listen")
	flag.Parse()

	db := models.NewInMemoryStorage()
	app.New(db).Start(*port)
}
