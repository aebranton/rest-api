package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aebranton/rest-api/internal/database"
	transHTTP "github.com/aebranton/rest-api/internal/transport/http"
	"github.com/aebranton/rest-api/internal/user"
)

// Server Constants (all in seconds)
const IdleTimeout = 120
const WriteTimeout = 1
const ReadTimeout = 1
const ShutdownGracePeriod = ReadTimeout * 3

// App - will contain things like our database connection
type App struct {
}

// Run - initializes application
func (app *App) Run() error {

	// Setup some basic logging
	l := log.New(os.Stdout, "rest-api ", log.LstdFlags)
	l.Println("App setup")

	// Create our database connection using our database package
	db, err := database.NewDatabase()
	if err != nil {
		return err
	}

	// Make sure we run our migrate function
	// currently only migrating users model, as this is all we have
	err = database.MigrateDB(db)
	if err != nil {
		return err
	}

	// Create a new user service that holds all our calls to manage users in the database
	// and supply our db pointer
	userService := user.NewService(db)

	// Creates our handler from our transport package.
	// The handler will contain a Router (gorillamux router) and needs a pointer to
	// our users service
	handler := transHTTP.NewHandler(userService)
	// Setup the rotues!
	handler.InitRoutes()

	// Tweak some paramters to make sure our connections dont get hung up for nonsense
	// We have very small data to read/write so timeouts are quite small
	server := http.Server{
		Addr:         ":8080",
		Handler:      handler.Router,
		IdleTimeout:  IdleTimeout * time.Second,
		ReadTimeout:  ReadTimeout * time.Second,
		WriteTimeout: WriteTimeout * time.Second,
	}

	// go func on this so we can do a graceful shutdown
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// Prepare a chan for os signals, so we can catch when something ie Ctrl+C is done on the service
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	// Waits for one of the signals to get pushed to this channel
	// Once a command is pushed to our signalChannel (ie ctrl+c) this will stop blocking, and
	// execute the code below
	sig := <-signalChannel

	// Notify the console we have recieved the kill sig and will just wait for the connections to close
	l.Println("Received kill signal - gracefully shutting down service via sig: ", sig)

	// Allows our erver to shutdown gracefully. In this case, we have a timeout of our Read timeout * 3
	// Once called, no new connections will be allowed, and existing connections will be allowed to finish their work
	// before we shutdown the service. Since its such a small read timeout, I did * 3 for some added buffer, though it shouldnt be necessary.
	killContext, _ := context.WithTimeout(context.Background(), ShutdownGracePeriod*time.Second)
	server.Shutdown(killContext)

	return nil
}

func main() {
	fmt.Println("Go REST API")

	// Setup our app (separated into a struct for easier testing and such later on)
	app := App{}
	err := app.Run()

	// Report any errors after run is complete
	if err != nil {
		fmt.Println("Error starting REST API")
		fmt.Println(err)
	}
}
