package main

import (
	"fmt"
	"net/http"
	"os"

	"app/config"
	controller "app/http/controllers"
	repository "app/repositories"
	route "app/routes"
	"app/services"

	"github.com/go-chi/chi/v5"
)

func main() {
	db := config.MustOpenDB(config.LoadDBConfig())
	defer db.Close()

	if err := config.EnsureSchema(db); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controller.NewAuthController(authService)

	postRepo := repository.NewPostRepository(db)
	postService := services.NewPostService(postRepo)
	postController := controller.NewPostController(postService)

	commentRepo := repository.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo, postRepo)
	commentController := controller.NewCommentController(commentService)

	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
	router.Mount("/", route.AuthRoutes(authController))
	router.Mount("/posts", route.PostRoutes(postController, commentController))

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
