package bootstrap

import (
	"fmt"
	"net/http"

	"github.com/asiman161/re-platform/app/replatform"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func SetupServer(i *replatform.Implementation, cfg AppConfig) *http.Server {
	r := chi.NewRouter()

	registerMiddlewares(r)
	registerHandlers(r, i)

	server := http.Server{Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), Handler: r}
	return &server
}

func registerMiddlewares(r *chi.Mux) {
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "First-name", "Last-name", "Email", "User-id"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Logger)
}

func registerHandlers(r *chi.Mux, i *replatform.Implementation) {
	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", i.Ping)
		r.Get("/users", i.Users)

		r.Route("/rooms", func(r chi.Router) {
			r.Get("/", i.GetRooms)
			r.Post("/", i.CreateRoom)
			r.Route("/{room_ID}", func(r chi.Router) {
				r.Get("/", i.GetRoom)
				r.HandleFunc("/ws", i.Room)
				r.Post("/close", i.CloseRoom)
				r.Get("/chat", i.GetMessages)
				r.Route("/pools", func(r chi.Router) {
					r.Post("/", i.CreateRoom)
					r.Get("/", i.GetPools)
					r.Post("/{pool_ID}/answer", i.AnswerPool)
					r.Post("/{pool_ID}/close", i.ClosePool)
				})
			})
		})

	})
}
