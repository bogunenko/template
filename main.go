package main

import (
	"log"
	"os"

	"github.com/bogunenko/template/engine"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	engine, err := engine.NewEngine(os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatal("cannot create engine:", err)
	}

	err = engine.Start()
	if err != nil {
		log.Fatal("cannot start engine:", err)
	}
}
