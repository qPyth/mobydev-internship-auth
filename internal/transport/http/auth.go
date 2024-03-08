package http

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) InitAuthRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", h.SignUp)
		r.Post("/sign-in", h.SignIn)
	})
}

type userSignUpReq struct {
	Email    string `json:"email" binding:"required, max=64"`
	Password string `json:"password" binding:"required, min=8,max=64"`
	PassConf string `json:"pass_conf" binding:"required, eqfield=Password"`
}

type userSignInReq struct {
	Email    string `json:"email" binding:"required, max=64"`
	Password string `json:"password" binding:"required, min=8,max=64"`
}

type SignInResp struct {
	Token string `json:"token"`
}

type AuthService interface {
	SignUp(ctx context.Context, email, password string) error
	SignIn(ctx context.Context, email, password string) (token string, err error)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req userSignUpReq
	if err := h.bindData(r, &req); err != nil {
		//TODO:handle error
		h.error(w, http.StatusBadRequest, err)
		return
	}

	if err := h.authService.SignUp(ctx, req.Email, req.Password); err != nil {
		//TODO:handle error
		h.error(w, http.StatusInternalServerError, err)
		return
	}
	_, err := w.Write([]byte("ok"))
	if err != nil {
		h.error(w, http.StatusInternalServerError, internalSrvErrorMsg)
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req userSignInReq
	if err := h.bindData(r, &req); err != nil {
		//TODO:handle error
		h.error(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.authService.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		//TODO:handle error
		h.error(w, http.StatusInternalServerError, err)
		return
	}
	h.NewResponse(w, http.StatusOK, SignInResp{Token: token})
}
