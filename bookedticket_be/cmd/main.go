package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"bookedticket_backend/internal/config"
	"bookedticket_backend/internal/handler"
	"bookedticket_backend/internal/repository"
	"bookedticket_backend/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	// เชื่อมต่อ Database
	db, err := repository.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	ticketRepo := repository.NewTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepo)
	ticketHandler := handler.NewTicketHandler(ticketService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // ระบุ origin ที่อนุญาต
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", ticketHandler.HealthCheck)

	authRequired := r.Group("/api/v1")
	{
		authRequired.GET("/ticket", ticketHandler.GetTicket)
		authRequired.GET("/ticket/:id", ticketHandler.GetTicketByID)
		authRequired.POST("/ticket", ticketHandler.CreateTicket)
		authRequired.PUT("/ticket/:id", ticketHandler.UpdateTicket)
	}

	port := cfg.APIPORT
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
