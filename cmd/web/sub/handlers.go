package sub

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/bejaneps/trading212/internal/crud"
	"github.com/bejaneps/trading212/internal/models"
)

func (e *env) handleStockPOST(w http.ResponseWriter, r *http.Request) {
	var err error

	comm := &crud.Commodity{}

	// Unmarshal json request from client
	err = json.NewDecoder(r.Body).Decode(comm)
	if err != nil {
		// log error
		log.Errorln(err)

		resp := &response{
			Success: false,
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		}

		resp.do(w)
		return
	}

	// Navigate back to trading page when request is done
	defer e.wd.Navigate(models.DemoTradingURL)

	var output string
	// Buy order
	if strings.ToLower(comm.Order) == "buy" {
		output, err = comm.Buy(e.wd.WebDriver)
		if err != nil {
			resp := &response{
				Success: false,
				Error:   err.Error(),
				Status:  http.StatusInternalServerError,
			}

			resp.do(w)
			return
		}
	}

	// Sell order
	if strings.ToLower(comm.Order) == "sell" {
		// TODO:
	}

	// Send successful response to client
	resp := &response{
		Success: true,
		Error:   "",
		Status:  http.StatusOK,
		Data:    output,
	}

	resp.do(w)
	return
}
