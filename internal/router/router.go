package router

import (
	"pillar/internal/handlers"
	"pillar/internal/handlers/auth"
	"pillar/internal/handlers/auth/admin"
	"pillar/internal/handlers/auth/normal_user"
	"pillar/internal/handlers/books"
	"pillar/internal/handlers/dashboard"
	"pillar/internal/handlers/images"
	"pillar/internal/handlers/lends"
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
	new_router.GET("/books/:id", books.SendBookById)
	new_router.GET("/books/:id/reviews", books.SendReviewsById)
	new_router.GET("/books", auth.TokenValidationMiddleware(books.SendPopularBooks))
	new_router.POST("/books/reviews", auth.TokenValidationMiddleware(books.AddReviews))
	new_router.POST("/orders", auth.TokenValidationMiddleware(books.RegisterOrder))

	// LISTS BOOK
	new_router.GET("/books/:id/lists/", auth.TokenValidationMiddleware(lists.GetUserLists))
	new_router.POST("/lists/:list_id/books", auth.TokenValidationMiddleware(lists.AddDocToList))
	new_router.POST("/lists", auth.TokenValidationMiddleware(lists.CreateList))
	new_router.DELETE("/lists/:list_id", auth.TokenValidationMiddleware(lists.EliminateList))
	new_router.PUT("/lists/:list_id/books", auth.TokenValidationMiddleware(lists.RenameList))
	new_router.DELETE("/lists/:list_id/books/:book_id", auth.TokenValidationMiddleware(lists.DeleteDocFromList))

	// LISTS SCREEN
	new_router.GET("/lists", auth.TokenValidationMiddleware(lists.GetAllLists))

	// LIKES
	new_router.PUT("/like/review/:id", auth.TokenValidationMiddleware(likes.AddLikeReview))
	new_router.DELETE("/like/review/:id", auth.TokenValidationMiddleware(likes.RemoveLikeReview))
	new_router.PUT("/like/list/:id", auth.TokenValidationMiddleware(likes.AddLikeList))
	new_router.DELETE("/like/list/:id", auth.TokenValidationMiddleware(likes.RemoveLikeList))

	// DASHBOARD
	new_router.GET("/dashboard/languages", auth.TokenValidationMiddleware(dashboard.GetLanguages))
	new_router.GET("/dashboard/publishers", auth.TokenValidationMiddleware(dashboard.GetPublishers))
	new_router.GET("/dashboard/countries", auth.TokenValidationMiddleware(dashboard.GetCountries))
	new_router.GET("/dashboard/formats", auth.TokenValidationMiddleware(dashboard.GetFormats))
	new_router.GET("/dashboard/authors", auth.TokenValidationMiddleware(dashboard.GetAuhors))
	new_router.GET("/dashboard/genders", auth.TokenValidationMiddleware(dashboard.GetGenders))
	new_router.GET("/dashboard/documents", auth.TokenValidationMiddleware(dashboard.GetDocuments))
	new_router.GET("/dashboard/orders", auth.TokenValidationMiddleware(dashboard.GetOrders))
	new_router.GET("/dashboard/authors/:document_id", auth.TokenValidationMiddleware(dashboard.GetAuthorsByDoc))

	new_router.POST("/dashboard/document", dashboard.AddDocToDB)
	new_router.POST("/dashboard/author", auth.TokenValidationMiddleware(dashboard.AddAuthor))
	new_router.POST("/dashboard/publisher", auth.TokenValidationMiddleware(dashboard.AddPublisher))
	new_router.POST("/dashboard/author/:author_id/document/:document_id", auth.TokenValidationMiddleware(dashboard.AddAuthorDocument))
	new_router.POST("/dashboard/order/:order_id", auth.TokenValidationMiddleware(dashboard.FinishOrder))

	// SEARCH
	new_router.GET("/search/:lookup", search.Search)

	// IMAGE
	new_router.GET("/cover/:document_id", images.ImageLink)

	// LENDS
	new_router.GET("/lends", auth.TokenValidationMiddleware(lends.GetLendsByUser))

	return new_router
}
