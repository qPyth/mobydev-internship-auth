package http

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/qPyth/mobydev-internship-auth/internal/domain"
	"github.com/qPyth/mobydev-internship-auth/internal/validators"
	"net/http"
)

func (h *Handler) InitUserRoutes(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/signup", h.SignUp)
		r.Post("/signin", h.SignIn)
		r.Post("/profile/update", h.UserProfileUpdate)
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

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)

type UserService interface {
	SignUp(ctx context.Context, email, password string) error
	SignIn(ctx context.Context, email, password string) (token string, err error)
	UpdateUserProfile(ctx context.Context, req domain.UserProfileUpdateReq) error
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req userSignUpReq
	if err := h.bindData(r, &req); err != nil {
		//TODO:handle error
		h.error(w, http.StatusBadRequest, err)
		return
	}
	if err := userSignUpReqValidation(req); err != nil {
		if errors.Is(err, ErrInvalidEmail) || errors.Is(err, ErrInvalidPassword) {
			h.error(w, http.StatusBadRequest, err)
			return
		}
		h.error(w, http.StatusInternalServerError, internalSrvErrorMsg)
		return
	}
	if err := h.userService.SignUp(ctx, req.Email, req.Password); err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			h.error(w, http.StatusBadRequest, err)
			return
		}
		h.error(w, http.StatusInternalServerError, internalSrvErrorMsg)
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
		h.error(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.userService.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			h.error(w, http.StatusBadRequest, err)
			return
		}
		h.error(w, http.StatusInternalServerError, err)
		return
	}
	h.NewResponse(w, http.StatusOK, SignInResp{Token: token})
}

func (h *Handler) UserProfileUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req domain.UserProfileUpdateReq
	if err := h.bindData(r, &req); err != nil {
		h.error(w, http.StatusBadRequest, err)
		return
	}
	err := h.userService.UpdateUserProfile(ctx, req)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			h.error(w, http.StatusBadRequest, err)
			return
		}
		h.error(w, http.StatusInternalServerError, internalSrvErrorMsg)
		return
	}
}

func userSignUpReqValidation(req userSignUpReq) error {
	emailValid, err := validators.EmailIsValid(req.Email)
	if err != nil {
		return err
	}

	passValid, err := validators.PasswordIsValid(req.Password)
	if err != nil {
		return err
	}

	switch {
	case !emailValid:
		return ErrInvalidEmail
	case !passValid:
		return ErrInvalidPassword
	}

	return nil
}
