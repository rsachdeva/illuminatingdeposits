package interestsvc

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits/mid"
	"github.com/rsachdeva/illuminatingdeposits/platform/web"
)

// Register constructs an http.Handler with all application routes defined.
func Register(shutdown chan os.Signal, db *sqlx.DB, log *log.Logger) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, log, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	{
		// Register health check handler. This route is not authenticated.
		c := Check{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		// Register user interestsvc.
		u := Users{db: db}

		// The route can't be authenticated because we need this route to
		// create the user in the first place.
		app.Handle(http.MethodPost, "/v1/users", u.Create)
	}

	{
		i := Interest{log: log}
		app.Handle(http.MethodPost, "/v1/interests", i.Create)
	}

	return app
}
