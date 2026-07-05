package route

import (
	controller "app/http/controllers"
	middleware "app/http/middlewares"

	"github.com/go-chi/chi/v5"
)

func PostRoutes(postController *controller.PostController, commentController *controller.CommentController) chi.Router {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		r.Post("/", postController.Create)
		r.Put("/{id}", postController.Update)
		r.Delete("/{id}", postController.Delete)
	})

	router.Get("/", postController.GetAll)
	router.Get("/{id}", postController.GetByID)

	router.Post("/{postID}/comments", commentController.Create)
	router.Get("/{postID}/comments", commentController.GetAllByPost)

	return router
}
