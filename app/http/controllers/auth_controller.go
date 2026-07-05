package controller

import (
	"errors"
	"net/http"

	helper "app/helpers"
	authrequests "app/http/requests/auth"
	"app/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req authrequests.RegisterRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if fieldErrors := helper.ValidateStruct(req); fieldErrors != nil {
		helper.WriteResponse(w, http.StatusUnprocessableEntity, fieldErrors[0].Message, fieldErrors)
		return
	}

	user, err := c.authService.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			helper.WriteResponse(w, http.StatusConflict, "email already exists", nil)
			return
		}
		helper.WriteResponse(w, http.StatusInternalServerError, "failed to create user", nil)
		return
	}

	helper.WriteResponse(w, http.StatusCreated, "user created", user)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req authrequests.LoginRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if fieldErrors := helper.ValidateStruct(req); fieldErrors != nil {
		helper.WriteResponse(w, http.StatusUnprocessableEntity, fieldErrors[0].Message, fieldErrors)
		return
	}

	user, token, err := c.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			helper.WriteResponse(w, http.StatusUnauthorized, "invalid credentials", nil)
			return
		}
		helper.WriteResponse(w, http.StatusInternalServerError, "failed to login", nil)
		return
	}

	helper.WriteResponse(w, http.StatusOK, "login success", map[string]any{
		"user":  user,
		"token": token,
	})
}
