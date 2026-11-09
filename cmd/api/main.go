package main

import (
	"log"
	"os"

	"github.com/FerreirAndre/estoque-caridade/internal/database"
	"github.com/gin-gonic/gin"
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

	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	println("conectado ao banco com sucesso")

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "oi caralho",
		})
	})

	router.Run(":3000")
}
