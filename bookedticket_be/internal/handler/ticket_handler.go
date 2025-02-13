package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"bookedticket_backend/internal/config"
	"bookedticket_backend/internal/models"
	"bookedticket_backend/internal/repository"
	"bookedticket_backend/internal/service"
	"log"
)

type TicketHandler struct {
	ticketService service.TicketService
}

// HealthCheck handles the health check endpoint
func (h *TicketHandler) HealthCheck(c *gin.Context) {
	cfg := config.LoadConfig()
	// เชื่อมต่อ Database
	db, err := repository.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()
	if err := repository.CheckDBConnection(db); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"detail": "Database connection failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "database": "connected"})
}

func NewTicketHandler(us service.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: us}
}

func (th *TicketHandler) GetTicket(c *gin.Context) {
	tickets, err := th.ticketService.GetTicket()
	if err != nil {
		log.Printf("Error fetching ticket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

func (th *TicketHandler) GetTicketByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	ticket, err := th.ticketService.GetTicketByID(id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

func (th *TicketHandler) CreateTicket(c *gin.Context) {
	var t models.TicketInsert
	if err := c.ShouldBindJSON(&t); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//ตรวจสอบค่า status ก่อน INSERT
	validStatuses := map[string]bool{
		"pending":  true,
		"accepted": true,
		"resolved": true,
		"rejected": true,
	}
	if !validStatuses[t.Status] {
		log.Println("Invalid status:", t.Status)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	if err := th.ticketService.CreateTicket(&t); err != nil {
		log.Println("Error inserting ticket:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	log.Println("Ticket created successfully:", t)
	c.JSON(http.StatusCreated, gin.H{"message": "Ticket created"})
}

func (th *TicketHandler) UpdateTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Invalid ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var ut models.Ticket
	if err := c.ShouldBindJSON(&ut); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validStatuses := map[string]bool{
		"pending":  true,
		"accepted": true,
		"resolved": true,
		"rejected": true,
	}
	if !validStatuses[ut.Status] {
		log.Println("Invalid status:", ut.Status)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	ticketUpdate := models.TicketInsert{
		Title:       ut.Title,
		Description: ut.Description,
		Contact:     ut.Contact,
		Status:      ut.Status,
	}

	if err := th.ticketService.UpdateTicket(id, ticketUpdate); err != nil {
		log.Println("Error updating ticket:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	log.Println("Ticket updated successfully:", ticketUpdate)
	c.JSON(http.StatusOK, gin.H{"message": "Ticket updated"})
}
