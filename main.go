package main

import (
	"github.com/57ajay/goTask/cmd"
	"github.com/57ajay/goTask/db"
	"log"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	cmd.Execute()
}
