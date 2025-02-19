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
	base        BaseHandler
	userUsecase usecase.UserUsecase
	jwtService  service.JWTService
}

func NewUserHandler(u usecase.UserUsecase, j service.JWTService) *UserHandler {
	return &UserHandler{
		base:        BaseHandler{},
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
		h.base.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.userUsecase.Register(input.Name, input.Email, input.Password); err != nil {
		h.base.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	h.base.ResponseJSON(c, http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.base.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.userUsecase.Login(input.Email, input.Password)
	if err != nil {
		h.base.ResponseUnauthorized(c, constants.InvalidCredentials)
		return
	}

	// Convert user.ID (string) to uint
	userID, err := strconv.ParseUint(user.ID, 10, 64)
	if err != nil {
		h.base.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	token, err := h.jwtService.GenerateToken(uint(userID))
	if err != nil {
		h.base.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	h.base.ResponseJSON(c, http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.base.ResponseError(c, http.StatusBadRequest, nil)
		return
	}

	user, err := h.userUsecase.GetUserByID(id)
	if err != nil {
		h.base.ResponseError(c, http.StatusNotFound, err)
		return
	}

	h.base.ResponseJSON(c, http.StatusOK, user)
}
