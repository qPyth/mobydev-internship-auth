package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/qPyth/mobydev-internship-auth/internal/domain"
	"log/slog"
	"net/http"
)

type UserService interface {
	UpdateUserProfile(req domain.UserProfileUpdateReq) error
}

type Handler struct {
	log         *slog.Logger
	authService AuthService
	userService UserService
	//TODO: tokenManager
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

var internalSrvErrorMsg = errors.New("server error")

func NewHandler(log *slog.Logger, authService AuthService, userService UserService) *Handler {
	return &Handler{log: log, authService: authService, userService: userService}
}

func (h *Handler) Init() *chi.Mux {
	r := chi.NewRouter()
	//TODO: add mws
	r.Use(middleware.Logger)
	h.InitAuthRoutes(r)
	h.InitUserRoutes(r)
	return r
}

func (h *Handler) bindData(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}
	return nil
}

func (h *Handler) NewResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("failed to encode response: ", "error", err.Error())
	}
}

func (h *Handler) error(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(status),
		Message: err.Error(),
	}); err != nil {
		h.log.Error("failed to encode error response: ", "error", err.Error())
	}

}
