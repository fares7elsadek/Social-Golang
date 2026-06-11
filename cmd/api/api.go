package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fares7elsadek/Social-Golang/internal/handler"
	"github.com/fares7elsadek/Social-Golang/internal/repository/postgres"
	service "github.com/fares7elsadek/Social-Golang/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func(app *application) mount(db *pgxpool.Pool) http.Handler {
	
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr) 
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/health", app.healthCheckHandler)

	userRepopsitory := postgres.NewUserRepository(db)
	postRepository := postgres.NewPostRepository(db)
	commentRepository := postgres.NewCommentRepository(db)

	userService := service.NewUserService(userRepopsitory)
	postService := service.NewPostService(postRepository,userRepopsitory)
	commentService := service.NewCommentService(commentRepository,userRepopsitory)

	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	r.Route("/api/v1",func (r chi.Router){
		// Users
		r.Post("/users", userHandler.CreateUser)
		r.Get("/users/{id}", userHandler.GetUserByID)
		r.Get("/users/email",userHandler.GetUserByEmail)
		r.Put("/users/{id}", userHandler.UpdateUser)
		r.Delete("/users/{id}", userHandler.DeleteUser)

		// Posts
		r.Post("/posts", postHandler.CreatePost)
		r.Get("/posts/{postId}", postHandler.GetPostByID)
		r.Get("/posts/{authorId}/author", postHandler.GetPostsByAuthorId)
		r.Put("/posts/{postId}", postHandler.UpdatePost)
		r.Delete("/posts/{postId}", postHandler.DeletePost)


		// Comments
		r.Post("/comments/{postId}/{authorId}", commentHandler.CreateComment)
		r.Get("/comments/{commentId}", commentHandler.GetCommentByID)
		r.Get("/comments/{postId}/post", commentHandler.GetCommentsByPostId)
		r.Put("/comments/{commentId}", commentHandler.UpdateComment)
		r.Delete("/comments/{commentId}", commentHandler.DeleteComment)
	})

	return r
}

func (app *application) run(handler http.Handler) error {

	srv := &http.Server{
		Addr: app.config.addr,
		Handler: handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}

	log.Printf("Server running on port %v",app.config.addr)

	return srv.ListenAndServe()
}