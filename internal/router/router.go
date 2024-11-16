package router

import (
	"pillar/internal/handlers"
	"pillar/internal/handlers/auth"
	"pillar/internal/handlers/auth/admin"
	"pillar/internal/handlers/auth/normal_user"
	"pillar/internal/handlers/books"
	"pillar/internal/handlers/dashboard"
	"pillar/internal/handlers/likes"
	"pillar/internal/handlers/lists"
	"pillar/internal/handlers/search"

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
	new_router.GET("/books/:id", auth.TokenValidationMiddleware(books.SendBookById))
	new_router.GET("/books/:id/reviews", auth.TokenValidationMiddleware(books.SendReviewsById))
	new_router.GET("/books", auth.TokenValidationMiddleware(books.SendPopularBooks))
	new_router.POST("/books/reviews", auth.TokenValidationMiddleware(books.AddReviews))
	new_router.POST("/orders", auth.TokenValidationMiddleware(books.RegisterOrder))

	// LISTS BOOK
	new_router.GET("/books/:id/lists/", auth.TokenValidationMiddleware(lists.GetUserLists))
	new_router.POST("/lists/:list_id/books", auth.TokenValidationMiddleware(lists.AddDocToList))
	new_router.POST("/lists", auth.TokenValidationMiddleware(lists.CreateList))
	new_router.PATCH("/lists/:list_id/books", auth.TokenValidationMiddleware(lists.RenameList))
	new_router.DELETE("/lists/:list_id/books/:book_id", auth.TokenValidationMiddleware(lists.DeleteDocFromList))

	// LISTS SCREEN
	new_router.GET("/lists", auth.TokenValidationMiddleware(lists.GetAllLists))

	// LIKES
	new_router.PUT("/like/review/:id", auth.TokenValidationMiddleware(likes.AddLikeReview))
	new_router.DELETE("/like/review/:id", auth.TokenValidationMiddleware(likes.RemoveLikeReview))
	new_router.PUT("/like/list/:id", auth.TokenValidationMiddleware(likes.AddLikeList))
	new_router.DELETE("/like/list/:id", auth.TokenValidationMiddleware(likes.RemoveLikeList))

	// DASHBOARD
	new_router.GET("/languages", auth.TokenValidationMiddleware(dashboard.GetLanguages))
	new_router.GET("/publishers", auth.TokenValidationMiddleware(dashboard.GetPublishers))
	new_router.POST("/document", auth.TokenValidationMiddleware(dashboard.AddDocToDB))

	// SEARCH
	new_router.GET("/search/:lookup", search.Search)

	return new_router
}
