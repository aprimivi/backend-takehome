package route

import (
	controller "app/http/controllers"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(authController *controller.AuthController) chi.Router {
	router := chi.NewRouter()

	router.Post("/register", authController.Register)
	router.Post("/login", authController.Login)

	return router
}
