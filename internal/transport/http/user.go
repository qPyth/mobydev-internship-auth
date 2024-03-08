package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) InitUserRoutes(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/profile/update", h.UserProfileUpdate)
	})
}

func (h *Handler) UserProfileUpdate(w http.ResponseWriter, r *http.Request) {

}
