package main

import (
	"context"
	"log"
	"os"

	"github.com/FerreirAndre/estoque-caridade/internal/database"
	"github.com/FerreirAndre/estoque-caridade/internal/handler"
	"github.com/FerreirAndre/estoque-caridade/internal/repository"
	"github.com/FerreirAndre/estoque-caridade/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("arquivo .env nao encontrado; utilizando variaveis de ambiente")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DatabaseURL nao encontrada")
	}
	if err := database.RunMigrations(databaseURL); err != nil {
		log.Fatal(err)
	}
	println("migrations executadas com sucesso")

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

	allowedOrigin := os.Getenv("FRONTEND_URL")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:4200"
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			allowedOrigin,
		},
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	handler.RegisterRoutes(router, itemService, movimentacaoService)

	if err := router.Run(os.Getenv("API_PORT")); err != nil {
		log.Fatal("erro ao iniciar api: ", err)
	}
}
