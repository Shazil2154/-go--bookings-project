package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Shazil2154/-go--bookings-project/internal/config"
	"github.com/Shazil2154/-go--bookings-project/internal/driver"
	"github.com/Shazil2154/-go--bookings-project/internal/handlers"
	"github.com/Shazil2154/-go--bookings-project/internal/helpers"
	"github.com/Shazil2154/-go--bookings-project/internal/models"
	"github.com/Shazil2154/-go--bookings-project/internal/render"
	"github.com/alexedwards/scs/v2"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}

// TODO:
// TO My Future and less exausted self
// I don't really like this pattern of initializing everything on the main function
// Maybe create a pkg to do the initialization or just seperate out this run function in the main pkg
// This is not js so we can pass actual references to mutate thing at the app level My suggestion is to add different files in the
// Main package like Template.go, Database.go, Config.go, ...e.t.c. or maybe just something Like Configurations.go with an init method
func run() (*driver.DB, error) {
	// What I am going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// Change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// Connect to database
	log.Println("Connecting to Database.....")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=Shazil2154")

	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	log.Println("Connected to the Database!")

	defer db.SQL.Close()

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Can not create template cache", err)
		return nil, err
	}

	app.TemplateCache = tc

	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
