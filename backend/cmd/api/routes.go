package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server api routes
func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFoundResponse(w, r)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.methodNotAllowedResponse(w, r)
	})

	r.Get("/api/v1/healthcheck", app.healthCheckHandler)

	r.Post("/api/v1/users/register", app.registerHandler)
	r.Post("/api/v1/users/login", app.loginHandler)

	return r
}

// Health check handler
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, envelope{"status": "OK"}, nil)
}
