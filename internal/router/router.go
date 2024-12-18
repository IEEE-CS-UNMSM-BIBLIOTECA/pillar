package router

import (
	"pillar/internal/handlers"
	"pillar/internal/handlers/auth"
	"pillar/internal/handlers/auth/admin"
	"pillar/internal/handlers/auth/normal_user"
	"pillar/internal/handlers/books"
	"pillar/internal/handlers/dashboard"
	"pillar/internal/handlers/dashboard/edit"
	"pillar/internal/handlers/images"
	"pillar/internal/handlers/lends"
	"pillar/internal/handlers/likes"
	"pillar/internal/handlers/lists"
	"pillar/internal/handlers/reviews"
	"pillar/internal/handlers/search"
	"pillar/internal/handlers/user"

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
	new_router.GET("/books", books.SendPopularBooks)
	new_router.POST("/books/reviews", auth.TokenValidationMiddleware(books.AddReviews))
	new_router.POST("/orders", auth.TokenValidationMiddleware(books.RegisterOrder))

	// REVIEWS
	new_router.GET("/reviews/:id", auth.TokenValidationMiddleware(reviews.ReviewByID))

	// LISTS BOOK
	new_router.GET("/books/:id/lists/", auth.TokenValidationMiddleware(lists.GetUserLists))
	new_router.POST("/lists/:list_id/books", auth.TokenValidationMiddleware(lists.AddDocToList))
	new_router.POST("/lists", auth.TokenValidationMiddleware(lists.CreateList))
	new_router.DELETE("/lists/:list_id", auth.TokenValidationMiddleware(lists.EliminateList))
	new_router.PUT("/lists/:list_id/books", auth.TokenValidationMiddleware(lists.RenameList))
	new_router.DELETE("/lists/:list_id/books/:book_id", auth.TokenValidationMiddleware(lists.DeleteDocFromList))

	// LISTS SCREEN
	new_router.GET("/lists", auth.TokenValidationMiddleware(lists.GetAllLists))
	new_router.GET("/lists/:list_id/books", books.GetBooksFromList)
	new_router.GET("/list/:list_id", auth.TokenValidationMiddleware(lists.GetListById))

	// LIKES
	new_router.POST("/reviews/:id/like", auth.TokenValidationMiddleware(likes.AddLikeReview))
	new_router.DELETE("/reviews/:id/like", auth.TokenValidationMiddleware(likes.RemoveLikeReview))
	new_router.POST("/lists/:list_id/like", auth.TokenValidationMiddleware(likes.AddLikeList))
	new_router.DELETE("/lists/:list_id/like", auth.TokenValidationMiddleware(likes.RemoveLikeList))

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
	new_router.GET("/dashboard/alltags/", auth.TokenValidationMiddleware(dashboard.GetAllTags))

	new_router.POST("/dashboard/document", auth.TokenValidationMiddleware(dashboard.AddDocToDB))
	new_router.POST("/dashboard/author", auth.TokenValidationMiddleware(dashboard.AddAuthor))
	new_router.POST("/dashboard/publisher", auth.TokenValidationMiddleware(dashboard.AddPublisher))
	new_router.POST("/dashboard/author/:author_id/document/:document_id", auth.TokenValidationMiddleware(dashboard.AddAuthorDocument))
	new_router.POST("/dashboard/order/:order_id", auth.TokenValidationMiddleware(dashboard.FinishOrder))
	new_router.POST("/dashboard/tag", auth.TokenValidationMiddleware(dashboard.AddTag))

	new_router.PATCH("/dashboard/document/edit", auth.TokenValidationMiddleware(edit.EditDoc))

	// SEARCH
	new_router.GET("/search/:lookup", search.Search)

	// IMAGE
	new_router.GET("/cover/:document_id", images.ImageLink)

	// LENDS
	new_router.GET("/lends", auth.TokenValidationMiddleware(lends.GetLendsByUser))

	// USER LISTS AND REVIEWS
	new_router.GET("/user/:user_id/screen", auth.TokenValidationMiddleware(user.GetUserById))
	new_router.GET("/user/:user_id/reviews", auth.TokenValidationMiddleware(reviews.GetReviewsByUserId))
	new_router.GET("/user/:user_id/lists", auth.TokenValidationMiddleware(lists.GetListByUserId))

	return new_router
}
