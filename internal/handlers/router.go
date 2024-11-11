package handlers

import (
	"github.com/julienschmidt/httprouter"
)

func NewPillarRouter() *httprouter.Router {
	new_router := httprouter.New()
	new_router.GET("/login", HndLogin)
	new_router.POST("/register", HndSignUp)

	new_router.OPTIONS("/document/:field", HndOptGetDocumentsBy)
	new_router.POST("/document/:field", HndGetDocumentsBy)

	new_router.GET("/popular-books", sendPopularBooks)
	new_router.GET("/books/:id", sendBookById)

	return new_router
}
