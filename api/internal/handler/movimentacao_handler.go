package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/FerreirAndre/estoque-caridade/internal/model"
	"github.com/FerreirAndre/estoque-caridade/internal/service"
	"github.com/gin-gonic/gin"
)

type MovimentacaoHandler struct {
	service *service.MovimentacaoService
}

func NewMovimentacaoHandler(service *service.MovimentacaoService) *MovimentacaoHandler {
	return &MovimentacaoHandler{
		service: service,
	}
}

type MovimentacaoRequest struct {
	ItemID     int64                      `json:"item_id" binding:"required"`
	Natureza   model.NaturezaMovimentacao `json:"natureza" binding:"required"`
	Quantidade int64                      `json:"quantidade" binding:"required"`
}

func (h *MovimentacaoHandler) Create(c *gin.Context) {
	var req MovimentacaoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos",
		})
		return
	}

	movimentacao := model.Movimentacao{
		ItemID:     req.ItemID,
		Natureza:   req.Natureza,
		Quantidade: req.Quantidade,
	}

	err := h.service.Create(c.Request.Context(), &movimentacao)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})

		case errors.Is(err, service.ErrInvalidQuantity),
			errors.Is(err, service.ErrInvalidNature):
			c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})

		case errors.Is(err, service.ErrInsufficientStock):
			c.JSON(http.StatusConflict, gin.H{"erro": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"erro": "Erro ao registrar movimentação",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, movimentacao)
}

func (h *MovimentacaoHandler) FindAll(c *gin.Context) {
	movimentacoes, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao buscar movimentações",
		})
		return
	}

	c.JSON(http.StatusOK, movimentacoes)
}

func (h *MovimentacaoHandler) FindByItemId(c *gin.Context) {
	itemID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || itemID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido",
		})
		return
	}

	movimentacoes, err := h.service.FindByItemId(
		c.Request.Context(),
		itemID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao buscar movimentações",
		})
		return
	}

	c.JSON(http.StatusOK, movimentacoes)
}
