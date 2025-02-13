package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"bookedticket_backend/internal/config"
	"bookedticket_backend/internal/models"

	_ "github.com/lib/pq"
)

type TicketRepository interface {
	GetTicket() ([]models.Ticket, error)
	GetTicketByID(id int) (models.Ticket, error)
	CreateTicket(ticket *models.TicketInsert) error
	UpdateTicket(id int, ut models.TicketInsert) error
}

type ticketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func ConnectDB(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	return sql.Open("postgres", dsn)
}

func CheckDBConnection(db *sql.DB) error {
	return db.Ping()
}

func (ur *ticketRepository) GetTicket() ([]models.Ticket, error) {
	rows, err := ur.db.Query(`SELECT id, title, description, contact, status, created_ticket_at, lastest_updated_at FROM ticket ORDER BY lastest_updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets := []models.Ticket{}
	for rows.Next() {
		ticket := models.Ticket{}
		err := rows.Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Contact, &ticket.Status, &ticket.Created_ticket_at, &ticket.Lastest_updated_at)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

func (ur *ticketRepository) GetTicketByID(id int) (models.Ticket, error) {
	row := ur.db.QueryRow("SELECT id, title, description, contact, status, created_ticket_at, lastest_updated_at FROM ticket WHERE id = $1", id)
	ticket := models.Ticket{}
	err := row.Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Contact, &ticket.Status, &ticket.Created_ticket_at, &ticket.Lastest_updated_at)
	if err != nil {
		return models.Ticket{}, err
	}
	return ticket, nil
}

func (ur *ticketRepository) CreateTicket(ticket *models.TicketInsert) error {
	ticket.Status = strings.ToLower(ticket.Status)

	validStatuses := map[string]bool{
		"pending":  true,
		"accepted": true,
		"resolved": true,
		"rejected": true,
	}

	if !validStatuses[ticket.Status] {
		return fmt.Errorf("invalid status: %s", ticket.Status)
	}

	_, err := ur.db.Exec("INSERT INTO ticket (title, description, contact, status) VALUES ($1, $2, $3, $4)",
		ticket.Title, ticket.Description, ticket.Contact, ticket.Status)

	if err != nil {
		return err
	}
	return nil
}

func (ur *ticketRepository) UpdateTicket(id int, ut models.TicketInsert) error {
	_, err := ur.db.Exec("UPDATE ticket SET title = $1, description = $2, contact = $3, status = $4 WHERE id = $5", ut.Title, ut.Description, ut.Contact, ut.Status, id)
	if err != nil {
		return err
	}
	return nil
}
