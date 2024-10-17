package handlers

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)


func HndOptGetDocumentsBy(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    w.WriteHeader(http.StatusOK)
}
