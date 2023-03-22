package main

import (
	"flag"

	"github.com/charmbracelet/log"

	"github.com/wuhan005/go-template/internal/db"
	"github.com/wuhan005/go-template/internal/route"
)

func main() {
	port := flag.Int("port", 8080, "port to listen")
	flag.Parse()

	db, err := db.Init()
	if err != nil {
		log.Fatal("Failed to initialize database", "error", err)
	}

	f := route.New(db)
	f.Run(*port)
}
