package options_user

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HndOptSignup(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}
