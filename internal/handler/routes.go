package handler

import (
	"github.com/FerreirAndre/estoque-caridade/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	itemService *service.ItemService,
	movimentacaoService *service.MovimentacaoService,
) {
	itemHandler := NewItemHandler(itemService)
	movimentacaoHandler := NewMovimentacaoHandler(movimentacaoService)

	itens := router.Group("/itens")
	{
		itens.POST("", itemHandler.Create)
		itens.GET("", itemHandler.FindAll)
		itens.GET("/:id", itemHandler.FindById)
		itens.PUT("/:id", itemHandler.Update)
		itens.DELETE("/:id", itemHandler.Delete)
	}

	movimentacoes := router.Group("/movimentacoes")
	{
		movimentacoes.POST("", movimentacaoHandler.Create)
		movimentacoes.GET("", movimentacaoHandler.FindAll)
		movimentacoes.GET("/item/:id", movimentacaoHandler.FindByItemId)
	}
}
