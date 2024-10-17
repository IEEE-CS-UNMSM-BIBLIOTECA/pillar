package handlers

import "github.com/julienschmidt/httprouter"

func NewPillarRouter() *httprouter.Router {
    new_router := httprouter.New()
    new_router.OPTIONS("/document/:field", HndOptGetDocumentsBy)
    new_router.POST("/document/:field", HndGetDocumentsBy)

    return new_router
}



