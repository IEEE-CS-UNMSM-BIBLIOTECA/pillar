package router

import (
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

func secure(handler httprouter.Handle) httprouter.Handle {
	return auth.TokenValidationMiddleware(handler)
}

func registerAuthRoutes(r *httprouter.Router) {
	// AUTH USERS
	r.POST("/login", normal_user.HndLogin)
	r.POST("/register", normal_user.HndSignUp)

	// AUTH ADMIN
	r.POST("/login/admin", admin.HndLoginAdmin)

	// AUTH
	r.GET("/protected", secure(auth.HndProtectedEndpoint))
}

func registerBooksRoutes(r *httprouter.Router) {
	r.GET("/books/:id", secure(books.SendBookById))
	r.GET("/books/:id/reviews", secure(books.SendReviewsById))
	r.GET("/books", books.SendPopularBooks)
	r.POST("/books/reviews", secure(books.AddReviews))
	r.POST("/orders", secure(books.RegisterOrder))
}

func registerListsBooksRoutes(r *httprouter.Router) {
	// LISTS BOOK
	r.GET("/books/:id/lists/", secure(lists.GetUserLists))
	r.POST("/lists/:list_id/books", secure(lists.AddDocToList))
	r.POST("/lists", secure(lists.CreateList))
	r.DELETE("/lists/:list_id", secure(lists.EliminateList))
	r.PUT("/lists/:list_id/books", secure(lists.RenameList))
	r.DELETE("/lists/:list_id/books/:book_id", secure(lists.DeleteDocFromList))
}

func registerListsScreen(r *httprouter.Router) {
	r.GET("/lists", secure(lists.GetAllLists))
	r.GET("/lists/:list_id/books", books.GetBooksFromList)
	r.GET("/list/:list_id", secure(lists.GetListById))
}

func registerLikesRoutes(r *httprouter.Router) {
	r.POST("/reviews/:id/like", secure(likes.AddLikeReview))
	r.DELETE("/reviews/:id/like", secure(likes.RemoveLikeReview))
	r.POST("/lists/:list_id/like", secure(likes.AddLikeList))
	r.DELETE("/lists/:list_id/like", secure(likes.RemoveLikeList))
}

func registerDashboardRoutes(r *httprouter.Router) {
	r.GET("/dashboard/languages", secure(dashboard.GetLanguages))
	r.GET("/dashboard/publishers", secure(dashboard.GetPublishers))
	r.GET("/dashboard/countries", secure(dashboard.GetCountries))
	r.GET("/dashboard/formats", secure(dashboard.GetFormats))
	r.GET("/dashboard/authors", secure(dashboard.GetAuhors))
	r.GET("/dashboard/genders", secure(dashboard.GetGenders))
	r.GET("/dashboard/documents", secure(dashboard.GetDocuments))
	r.GET("/dashboard/orders", secure(dashboard.GetOrders))
	r.GET("/dashboard/authors/:document_id", secure(dashboard.GetAuthorsByDoc))
	r.GET("/dashboard/alltags/", secure(dashboard.GetAllTags))

	r.POST("/dashboard/document", secure(dashboard.AddDocToDB))
	r.POST("/dashboard/author", secure(dashboard.AddAuthor))
	r.POST("/dashboard/publisher", secure(dashboard.AddPublisher))
	r.POST("/dashboard/author/:author_id/document/:document_id", secure(dashboard.AddAuthorDocument))
	r.POST("/dashboard/order/:order_id", secure(dashboard.FinishOrder))
	r.POST("/dashboard/tag", secure(dashboard.AddTag))

	r.PATCH("/dashboard/document/edit", secure(edit.EditDoc))
}

func NewPillarRouter() *httprouter.Router {
	new_router := httprouter.New()

	registerAuthRoutes(new_router)
	registerBooksRoutes(new_router)
	registerListsBooksRoutes(new_router)
	registerListsScreen(new_router)
	registerLikesRoutes(new_router)
	registerDashboardRoutes(new_router)

	// REVIEWS
	new_router.GET("/reviews/:id", secure(reviews.ReviewByID))

	// SEARCH
	new_router.GET("/search/:lookup", search.Search)

	// IMAGE
	new_router.GET("/cover/:document_id", images.ImageLink)

	// LENDS
	new_router.GET("/lends", secure(lends.GetLendsByUser))

	// USER LISTS AND REVIEWS
	new_router.GET("/user/:user_id/screen", secure(user.GetUserById))
	new_router.GET("/user/:user_id/reviews", secure(reviews.GetReviewsByUserId))
	new_router.GET("/user/:user_id/lists", secure(lists.GetListByUserId))

	return new_router
}
