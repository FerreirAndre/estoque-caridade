package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/FerreirAndre/estoque-caridade/internal/model"
	"github.com/FerreirAndre/estoque-caridade/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type ItemHandler struct {
	service *service.ItemService
}

func NewItemHandler(service *service.ItemService) *ItemHandler {
	return &ItemHandler{
		service: service,
	}
}

type ItemRequest struct {
	Nome          string              `json:"nome" binding:"required"`
	Unidade       model.UnidadeMedida `json:"unidade_medida" binding:"required"`
	EstoqueMinimo int64               `json:"estoque_minimo"`
}

func (h *ItemHandler) Create(c *gin.Context) {
	var req ItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos",
		})
		return
	}

	item := model.Item{
		Nome:          req.Nome,
		Unidade:       req.Unidade,
		EstoqueMinimo: req.EstoqueMinimo,
	}

	if err := h.service.Create(c.Request.Context(), &item); err != nil {
		log.Printf("erro ao criar item: %w", err)

		if h.handleItemError(c, err) {
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"erro": "erro ao cadastrar item",
		})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) FindAll(c *gin.Context) {
	itens, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao buscar itens",
		})
		return
	}

	c.JSON(http.StatusOK, itens)
}

func (h *ItemHandler) FindById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido",
		})
		return
	}

	item, err := h.service.FindById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao buscar item",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido",
		})
		return
	}

	var req ItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos",
		})
		return
	}

	item := model.Item{
		ID:            id,
		Nome:          req.Nome,
		Unidade:       req.Unidade,
		EstoqueMinimo: req.EstoqueMinimo,
	}

	if err := h.service.Update(c.Request.Context(), &item); err != nil {
		log.Printf("erro ao atualizar item: %w", err)

		if h.handleItemError(c, err) {
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao atualizar item",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "ID inválido",
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		log.Printf("erro ao excluir item: %w", err)

		if h.handleItemError(c, err) {
			return
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "nao foi possivel excluir este item porque existem movimentacoes registradas",
			})
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao excluir item",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ItemHandler) handleItemError(c *gin.Context, err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" &&
			pgErr.ConstraintName == "uq_itens_nome" {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "Já existe um item cadastrado com esse nome.",
			})
			return true
		}
	}

	if errors.Is(err, service.ErrItemNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": err.Error(),
		})
		return true
	}

	if errors.Is(err, service.ErrInvalidName) ||
		errors.Is(err, service.ErrInvalidUnit) ||
		errors.Is(err, service.ErrInvalidMinimumStock) {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": err.Error(),
		})
		return true
	}

	return false
}
