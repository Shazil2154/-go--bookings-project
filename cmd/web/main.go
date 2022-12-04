package main

import (
	"encoding/gob"
	"github.com/Shazil2154/-go--bookings-project/internal/config"
	"github.com/Shazil2154/-go--bookings-project/internal/handlers"
	"github.com/Shazil2154/-go--bookings-project/internal/helpers"
	"github.com/Shazil2154/-go--bookings-project/internal/models"
	"github.com/Shazil2154/-go--bookings-project/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// What I am going to put in the session
	gob.Register(models.Reservation{})

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

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Can not create template cache", err)
		return err
	}

	app.TemplateCache = tc

	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return nil
}
