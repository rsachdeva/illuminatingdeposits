// Package postgreshealth provides postgress health status check service
package postgreshealth

import (
	"context"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits-rest/jsonfmt"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"github.com/rsachdeva/illuminatingdeposits-rest/postgreshealth/healthvalue"
	"go.opencensus.io/trace"
)

// service provides support for orchestration health checks.
type service struct {
	db *sqlx.DB

	// ADD OTHER STATE LIKE THE LOGGER IF NEEDED.
}

// Health validates the jsonfmt is healthy and ready to accept requests.
func (c service) Health(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	log.Println("in Health service")
	ctx, span := trace.StartSpan(ctx, "postgresconn.service.Health")
	defer span.End()

	var ht healthvalue.Postgres

	// service if the postgresconn is ready.
	if err := healthvalue.StatusCheck(ctx, c.db); err != nil {

		// If the postgresconn is not ready we will tell the cli and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		ht.Status = "Postgres Db Not Ready"
		return jsonfmt.Respond(ctx, w, ht, http.StatusInternalServerError)
	}

	ht.Status = "Postgres Db Ok"
	return jsonfmt.Respond(ctx, w, ht, http.StatusOK)
}

func RegisterSvc(db *sqlx.DB, rt *muxhttp.Router) {
	// Register health check handler. This jsonfmt is not authenticated.
	c := service{db: db}
	log.Println("registering /v1/health route")
	rt.Handle(http.MethodGet, "/v1/health", c.Health)
}
