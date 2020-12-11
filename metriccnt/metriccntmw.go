// Package errconv provides counter for metrics for all services
package metriccnt

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"runtime"

	"github.com/rsachdeva/illuminatingdeposits-rest/appmux"
	"go.opencensus.io/trace"
)

// m contains the global program counters for the application.
var m = struct {
	gr  *expvar.Int
	req *expvar.Int
	err *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
	err: expvar.NewInt("errors"),
}

// Metrics updates program counters.
func Middleware() appmux.Middleware {

	// This is the actual middlewarefunc function to be executed.
	f := func(before appmux.Handler) appmux.Handler {

		// Wrap this handler around the next one provided.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Printf("\t\tEntering metricscnt Middleware handler is %T\n", before)
			defer fmt.Printf("\t\tExiting metricscnt Middleware handler is %T\n", before)

			ctx, span := trace.StartSpan(ctx, "metriccnt.Middleware")
			defer span.End()

			err := before(ctx, w, r)

			// Increment the request counter.
			m.req.Add(1)

			// Update the count for the number of active goroutines every 100 requests.
			if m.req.Value()%100 == 0 {
				m.gr.Set(int64(runtime.NumGoroutine()))
			}

			// Increment the errors counter if an error occurred on this request.
			if err != nil {
				m.err.Add(1)
			}

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return f
}
