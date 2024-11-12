package options_books

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HndOptBooks(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}
