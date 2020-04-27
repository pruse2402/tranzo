package routes

import (
	"tranzo/src/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func NewRouter(h *handlers.Provider) *chi.Mux {

	r := chi.NewRouter()

	r.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger, // Log API request calls
	)

	r.Post("/importExcel", h.ImportExcelFile)
	r.Get("/ping", h.Ping)

	return r
}
