package handler

import (
	"context"
	"fmt"
	"moneyget/internal/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type InvestmentHandler struct {
	BaseHandler
	investmentUsecase InvestmentUsecase
}

type InvestmentUsecase interface {
	CreateInvestment(ctx context.Context, userID string, amount float64, currency string, investmentType string, strategy string) error
	GetInvestment(ctx context.Context, id string) (*domain.Investment, error)
}

func NewInvestmentHandler(iu InvestmentUsecase) *InvestmentHandler {
	return &InvestmentHandler{
		investmentUsecase: iu,
	}
}

type CreateInvestmentRequest struct {
	UserID   string  `json:"user_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
	Type     string  `json:"type" binding:"required"`
	Strategy string  `json:"strategy" binding:"required"`
}

func (h *InvestmentHandler) CreateInvestment(c *gin.Context) {
	ctx, cancel := h.NewContext(c, 10*time.Second)
	defer cancel()

	var req CreateInvestmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.investmentUsecase.CreateInvestment(ctx, req.UserID, req.Amount, req.Currency, req.Type, req.Strategy); err != nil {
		h.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	h.ResponseJSON(c, http.StatusCreated, gin.H{"message": "Investment created successfully"})
}

func (h *InvestmentHandler) GetInvestment(c *gin.Context) {
	ctx, cancel := h.NewContext(c, 5*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		h.ResponseError(c, http.StatusBadRequest, fmt.Errorf("id is required"))
		return
	}

	investment, err := h.investmentUsecase.GetInvestment(ctx, id)
	if err != nil {
		h.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	h.ResponseJSON(c, http.StatusOK, investment)
}
