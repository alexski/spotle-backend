package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"spotle-backend/handler"

	"github.com/joho/godotenv"
)

// Current version of application
var Version = "0.0.1"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	a := handler.App{}
	a.Initialize(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"))

	port := "8080"
	_, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("Defaulting port...")
	} else {
		port = os.Getenv("PORT")
	}

	srv := &http.Server{
		Handler:      a.Router,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("API Version: " + Version)
	fmt.Println("Running server on port: " + port)
	log.Fatal(srv.ListenAndServe())
}
