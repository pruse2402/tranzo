package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"tranzo/src/handlers"
	route "tranzo/src/routes"

	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
)

func main() {

	//Configuring logger for the app
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logger := log.New(os.Stdout, "Tranzo : ", log.LstdFlags|log.Lshortfile)

	// Creating master session for MongoDB
	logger.Println("Initializing mongodb session...")

	dbSession := connectLocalDB()
	defer dbSession.Close()

	logger.Println("Initializing provider...")
	provider := handlers.NewProvider(logger, dbSession)

	logger.Println("Initializing routes...")
	router := route.NewRouter(provider)

	// Setting up CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "AuthKey", "if-modified-since", "Access-Control-Allow-Origin"},
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":8089"),
		Handler: c.Handler(router),
	}

	// Graceful shut down of server
	graceful := make(chan os.Signal)

	go func() {
		<-graceful
		logger.Println("Shutting down server...")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Fatalf("Could not do graceful shutdown: %v\n", err)
		}
	}()

	logger.Println("Listening server on : 8089")
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Listen: %s\n", err)
	}

	logger.Println("Server gracefully stopped")

}

func connectLocalDB() *mgo.Session {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Printf("Error in dialing mongo server: %s", err.Error())
	}
	return session
}
