package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	// cfg is the global configuration for the server. It's read in at startup from
	// the config.json file and enviornment variables, see config.go for more info.
	cfg *config

	// When was the last alert sent out?
	// Use this value to avoid bombing alerts
	lastAlertSent *time.Time

	// application database connection
	appDB *sql.DB
)

func init() {
	var err error
	cfg, err = initConfig()
	if err != nil {
		// panic if the server is missing a vital configuration detail
		panic(fmt.Errorf("server configuration error: %s", err.Error()))
	}
}

func main() {
	connectToAppDb()

	// what a wonderful phrase :)
	go startCrawling()

	// initialize a router to handle requests
	r := httprouter.New()

	// home handler, wrapped in middlware func
	r.GET("/", middleware(HandleDomains))

	// serve static content from public directory
	r.ServeFiles("/css/*filepath", http.Dir("public/css"))
	r.ServeFiles("/js/*filepath", http.Dir("public/js"))

	// print notable config settings
	printConfigInfo()

	// fire it up!
	fmt.Println("starting server on port", cfg.Port)
	// start server wrapped in a call to panic b/c http.ListenAndServe will not
	// return unless there's an error
	panic(http.ListenAndServe(":"+cfg.Port, r))
}
