package handlers

import (
	"pillar/internal/handlers/auth/normal_user"
	"pillar/internal/handlers/books"

	"github.com/julienschmidt/httprouter"
)

func NewPillarRouter() *httprouter.Router {
	new_router := httprouter.New()
	new_router.GET("/login", normal_user.HndLogin)
	new_router.POST("/register", normal_user.HndSignUp)

	new_router.OPTIONS("/document/:field", HndOptGetDocumentsBy)
	new_router.POST("/document/:field", HndGetDocumentsBy)

	new_router.POST("/books", books.SendPopularBooks)
	new_router.GET("/book/:id", books.SendBookById)
	new_router.GET("/book/:id/reviews", books.SendReviewsById)

	return new_router
}
