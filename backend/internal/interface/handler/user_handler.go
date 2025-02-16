package handler

import (
	"moneyget/internal/domain/service"
	"moneyget/internal/usecase"
	"net/http"
	"strconv"

	"moneyget/internal/domain/constants"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	BaseHandler
	userUsecase usecase.UserUsecase
	jwtService  service.JWTService
}

func NewUserHandler(u usecase.UserUsecase, j service.JWTService) *UserHandler {
	return &UserHandler{
		userUsecase: u,
		jwtService:  j,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.userUsecase.Register(input.Name, input.Email, input.Password); err != nil {
		h.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	h.ResponseJSON(c, http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.userUsecase.Login(input.Email, input.Password)
	if err != nil {
		h.ResponseUnauthorized(c, constants.InvalidCredentials)
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID)
	if err != nil {
		h.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	h.ResponseJSON(c, http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.userUsecase.GetUserByID(uint(id))
	if err != nil {
		h.ResponseError(c, http.StatusNotFound, err)
		return
	}

	h.ResponseJSON(c, http.StatusOK, user)
}
