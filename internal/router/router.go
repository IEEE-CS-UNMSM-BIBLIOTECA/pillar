package router

import (
	"pillar/internal/handlers"
	"pillar/internal/handlers/auth"
	"pillar/internal/handlers/auth/admin"
	"pillar/internal/handlers/auth/normal_user"
	"pillar/internal/handlers/books"
	"pillar/internal/handlers/lists"

	"github.com/julienschmidt/httprouter"
)

func NewPillarRouter() *httprouter.Router {
	new_router := httprouter.New()

	// AUTH USERS
	new_router.POST("/login", normal_user.HndLogin)
	new_router.POST("/register", normal_user.HndSignUp)

	// AUTH ADMIN
	new_router.POST("/login/admin", admin.HndLoginAdmin)

	// AUTH
	new_router.GET("/protected", auth.TokenValidationMiddleware(auth.HndProtectedEndpoint))

	// RODRO
	new_router.OPTIONS("/document/:field", handlers.HndGetDocumentsBy)
	new_router.POST("/document/:field", handlers.HndOptGetDocumentsBy)

	// BOOKS
	new_router.GET("/books/:id", books.SendBookById)
	new_router.GET("/books/:id/reviews", books.SendReviewsById)
	new_router.POST("/books", books.SendPopularBooks)
	new_router.POST("/books/new-review", books.AddReviews)
	new_router.POST("/books/lend", books.RegisterLend)

	// LISTS
	new_router.GET("/books/:id/lists/:user_id", lists.GetUserLists)
	new_router.POST("/books/lists", lists.AddDocToList)
	new_router.POST("/books/lists/rename", lists.RenameList)
	new_router.DELETE("/books/lists", lists.DeleteDocFromList)

	return new_router
}
