package handlers

import (
	"io"
	"log"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/internal/invest"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/inout"
)

//Interest handler
type Interest struct {
	Log *log.Logger
}

// Create investment calculates for all banks, sent to the desired writer in JSON format
func (ih Interest) Create(w io.Writer, nibs invest.NewInterestBanks, executionTimes int) error {
	var ibs invest.InterestBanks
	var err error
	for j := 0; j < executionTimes; j++ {
		ibs, err = invest.Delta(nibs)
	}
	if err != nil {
		return errors.Wrap(err, "create calculating for invest.NewInterestBanks")
	}
	return inout.OutputJSON(w, ibs)
}
