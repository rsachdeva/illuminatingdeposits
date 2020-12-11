// Package interestcal provides interest calculation service
package interestcal

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/debug"
	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal/interestvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/responder"
	"go.opencensus.io/trace"
)

// Service handler
type Service struct {
	Log *log.Logger
}

// ListCalculations decodes the body of a request to create interest calculations. The full
// banks and deposit details with generated 30 days Service fields are sent back in the response.
func (*Service) ListCalculations(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "interestcal.Service.ListCalculations")
	defer span.End()

	debug.Dump(r)
	var nin interestvalue.NewInterest
	if err := responder.Decode(r, &nin); err != nil {
		return errors.Wrap(err, "decoding new interest calculation request with banks and deposits")
	}

	in, err := nin.CalculateDelta()
	if err != nil {
		return errors.Wrap(err, "creating new interest calculations")
	}

	return responder.Respond(ctx, w, &in, http.StatusCreated)
}

func RegisterSvc(log *log.Logger, m *responder.ServeMux) {
	i := Service{Log: log}
	m.Handle(http.MethodPost, "/v1/interests", i.ListCalculations)
}
