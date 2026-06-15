package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fares7elsadek/Social-Golang/internal/handler"
	authmiddleware "github.com/fares7elsadek/Social-Golang/internal/middlewares"
	"github.com/fares7elsadek/Social-Golang/internal/repository/postgres"
	service "github.com/fares7elsadek/Social-Golang/internal/services"
	"github.com/fares7elsadek/Social-Golang/internal/services/auth"
	"github.com/fares7elsadek/Social-Golang/internal/services/token"
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

	config := token.DefaultConfig([]byte("testsecret"),[]byte("testsecret"))

	userRepopsitory := postgres.NewUserRepository(db)
	postRepository := postgres.NewPostRepository(db)
	commentRepository := postgres.NewCommentRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepo(db)

	userService := service.NewUserService(userRepopsitory)
	postService := service.NewPostService(postRepository,userRepopsitory)
	commentService := service.NewCommentService(commentRepository,userRepopsitory)
	tokenService := token.New(config,refreshTokenRepo)
	authService := auth.NewAuthService(userRepopsitory,tokenService)

	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)
	authHandler := handler.NewAuthHandler(authService)

	r.Route("/api/v1",func (r chi.Router){
		// auth
		r.Post("/auth/register",authHandler.Register)
		r.Post("/auth/login",authHandler.Login)
		r.Post("/auth/refresh",authHandler.Refresh)
		r.Post("/auth/logout",authHandler.Logout)
		r.With(authmiddleware.Authenticate(tokenService)).Get("/me",authHandler.Me)

		// Users
		r.Get("/users/{id}", userHandler.GetUserByID)
		r.Get("/users/email",userHandler.GetUserByEmail)
		
		r.With(authmiddleware.Authenticate(tokenService),
		authmiddleware.RequireSelf(func(r *http.Request) string {
		return chi.URLParam(r, "id")
		})).Put("/users/{id}", userHandler.UpdateUser)
		
		r.With(authmiddleware.Authenticate(tokenService),
		authmiddleware.RequireSelf(func(r *http.Request) string {
		return chi.URLParam(r, "id")
		})).Delete("/users/{id}", userHandler.DeleteUser)


		// Posts
		r.With(authmiddleware.Authenticate(tokenService)).Post("/posts", postHandler.CreatePost)
		r.With(authmiddleware.Authenticate(tokenService)).Get("/posts/{postId}", postHandler.GetPostByID)
		r.With(authmiddleware.Authenticate(tokenService)).Get("/posts/{authorId}/author", postHandler.GetPostsByAuthorId)

		
		r.With(authmiddleware.Authenticate(tokenService)).Put("/posts/{postId}", postHandler.UpdatePost)
		r.With(authmiddleware.Authenticate(tokenService)).Delete("/posts/{postId}", postHandler.DeletePost)


		// Comments
		r.With(authmiddleware.Authenticate(tokenService)).Post("/comments/{postId}/{authorId}", commentHandler.CreateComment)
		r.With(authmiddleware.Authenticate(tokenService)).Get("/comments/{commentId}", commentHandler.GetCommentByID)
		r.With(authmiddleware.Authenticate(tokenService)).Get("/comments/{postId}/post", commentHandler.GetCommentsByPostId)
		r.With(authmiddleware.Authenticate(tokenService)).Put("/comments/{commentId}", commentHandler.UpdateComment)
		r.With(authmiddleware.Authenticate(tokenService)).Delete("/comments/{commentId}", commentHandler.DeleteComment)
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