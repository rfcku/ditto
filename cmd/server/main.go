package server

import (
	"context"
	"go-api/pkg/db"
	"log"
	"net/http"

	"go-api/cmd/server/authenticator"
	"go-api/cmd/server/router"

	"github.com/joho/godotenv"
	// gin-swagger middleware
	// swagger embed files
)



func Run() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}
	
	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := router.New(auth)

	client := db.ConnectDB()
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	log.Print("Server listening on http://localhost:8080/")
	if err := http.ListenAndServe("0.0.0.0:8080", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}

