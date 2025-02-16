package handler

import (
	"context"
	"moneyget/internal/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	BaseHandler
	portfolioUsecase PortfolioUsecase
}

type PortfolioUsecase interface {
	GetUserPortfolio(ctx context.Context, userID string) (*domain.Portfolio, error)
}

func NewPortfolioHandler(pu PortfolioUsecase) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioUsecase: pu,
	}
}

func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	ctx, cancel := h.NewContext(c, 5*time.Second)
	defer cancel()

	userID, exists := c.Get("userID")
	if !exists {
		h.ResponseUnauthorized(c, "user not authenticated")
		return
	}

	portfolio, err := h.portfolioUsecase.GetUserPortfolio(ctx, userID.(string))
	if err != nil {
		h.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	h.ResponseJSON(c, http.StatusOK, portfolio)
}
