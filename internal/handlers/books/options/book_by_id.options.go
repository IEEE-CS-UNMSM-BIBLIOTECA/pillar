package options

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HndOptBookById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}
