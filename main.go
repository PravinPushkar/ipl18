package main

import (
	"log"
	"net/http"

	"github.wdf.sap.corp/I334816/ipl18/backend"
)

func main() {
	log.Println("Starting server on port 3000...")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", backend.SetupAndGetRouter()))
	// scraper.Start()
}
