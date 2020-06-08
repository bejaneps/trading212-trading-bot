package sub

import "net/http"

func (e *env) routes() {
	// router for buying or selling new stock.
	e.router.Handle("/commodity", http.HandlerFunc(e.handleStockPOST)).Methods("POST")
}
