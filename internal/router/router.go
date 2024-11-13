package router

import (
	"pillar/internal/handlers"
	"pillar/internal/handlers/auth/normal_user"
	"pillar/internal/handlers/auth/normal_user/options_user"
	"pillar/internal/handlers/books"
	"pillar/internal/handlers/books/options_books"

	"github.com/julienschmidt/httprouter"
)

func NewPillarRouter() *httprouter.Router {
	new_router := httprouter.New()
	// AUTH USERS
	new_router.POST("/login", normal_user.HndLogin)
	new_router.OPTIONS("/login", options_user.HndOptSignin)
	new_router.POST("/register", normal_user.HndSignUp)
	new_router.OPTIONS("/register", options_user.HndOptSignup)

	// RODRO
	new_router.OPTIONS("/document/:field", handlers.HndGetDocumentsBy)
	new_router.POST("/document/:field", handlers.HndOptGetDocumentsBy)

	// BOOKS
	new_router.POST("/books", books.SendPopularBooks)
	new_router.OPTIONS("/books", options_books.HndOptBooks)
	new_router.GET("/book/:id", books.SendBookById)
	new_router.OPTIONS("/book/:id", options_books.HndOptBookById)
	new_router.GET("/book/:id/reviews", books.SendReviewsById)
	new_router.POST("/book/new-review", books.AddReviews)
	new_router.POST("/book/loan", books.RegisterLoan)

	return new_router
}
