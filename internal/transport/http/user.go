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
		r.With(h.JWTAuthMiddleware).Post("/profile/update", h.UserProfileUpdate)
	})
}

type userSignUpReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	PassConf string `json:"pass_conf" required:"true"`
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
	ErrInvalidPassword = errors.New("invalid password len")
	ErrPassConfirm     = errors.New("passwords do not match")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidBDay     = errors.New("invalid birthdate or date format not in RFC3339")
	ErrInvalidPhone    = errors.New("invalid phone number")
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
		h.log.Error("failed to bind user sign up request: ", "error", err.Error())
		//TODO:handle error
		h.error(w, http.StatusBadRequest, err)
		return
	}
	if err := userSignUpReqValidation(req); err != nil {
		h.log.Error("failed to validate user sign up request: ", "error", err.Error())
		if errors.Is(err, ErrInvalidEmail) || errors.Is(err, ErrInvalidPassword) || errors.Is(err, ErrPassConfirm) {
			h.error(w, http.StatusBadRequest, err)
			return
		}
		h.error(w, http.StatusInternalServerError, internalSrvErrorMsg)
		return
	}
	if err := h.userService.SignUp(ctx, req.Email, req.Password); err != nil {
		h.log.Error("failed to sign up user: ", "error", err.Error())
		if errors.Is(err, domain.ErrEmailExists) {
			h.error(w, http.StatusBadRequest, err)
			return
		}
		h.error(w, http.StatusInternalServerError, internalSrvErrorMsg)
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
		h.log.Error("failed to bind user profile update request: ", "error", err.Error())
		h.error(w, http.StatusBadRequest, err)
		return
	}
	err := updateProfileValidation(req)
	if err != nil {
		h.log.Error("failed to validate user profile update request: ", "error", err.Error())
		if errors.Is(err, ErrInvalidName) || errors.Is(err, ErrInvalidBDay) || errors.Is(err, ErrInvalidPhone) || errors.Is(err, ErrInvalidEmail) {
			h.error(w, http.StatusBadRequest, err)
			return
		}
		h.error(w, http.StatusInternalServerError, internalSrvErrorMsg)
		return
	}
	err = h.userService.UpdateUserProfile(ctx, req)
	if err != nil {
		h.log.Error("failed to update user profile: ", "error", err.Error())
		if errors.Is(err, domain.ErrUserNotFound) {
			h.error(w, http.StatusBadRequest, err)
			return
		}
		h.error(w, http.StatusInternalServerError, internalSrvErrorMsg)
		return
	}
	_, err = w.Write([]byte("ok"))
	if err != nil {
		h.log.Error("failed to write response: ", "error", err.Error())
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
	case !validators.PasswordsMatch(req.Password, req.PassConf):
		return ErrPassConfirm
	}

	return nil
}

func updateProfileValidation(req domain.UserProfileUpdateReq) error {
	if req.Name != nil {
		if len(*req.Name) > 64 {
			return ErrInvalidName
		}
	}

	if req.Email != nil {
		emailValid, err := validators.EmailIsValid(*req.Email)
		if err != nil {
			return err
		}
		if !emailValid {
			return ErrInvalidEmail
		}
	}

	if req.BDay != nil {
		if !validators.BDayValidation(*req.BDay) {
			return ErrInvalidBDay
		}
	}

	if req.PhoneNumber != nil {
		phoneValid, err := validators.PhoneE164Validation(*req.PhoneNumber)
		if err != nil {
			return err
		}
		if !phoneValid {
			return ErrInvalidPhone
		}
	}
	return nil
}
