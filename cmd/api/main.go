package main

import (
	"context"
	"log"
	"os"

	"github.com/FerreirAndre/estoque-caridade/internal/handler"
	"github.com/FerreirAndre/estoque-caridade/internal/repository"
	"github.com/FerreirAndre/estoque-caridade/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env:", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DatabaseURL nao encontrada")
	}
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	itemRepository := repository.NewItemRepository(pool)

	itemService := service.NewItemService(itemRepository)
	movimentacaoService := service.NewMovimentacaoService(pool)

	router := gin.Default()
	router.Use(gin.Logger())

	handler.RegisterRoutes(router, itemService, movimentacaoService)

	if err := router.Run(":3000"); err != nil {
		log.Fatal("erro rodando essa bosta")
	}
}
