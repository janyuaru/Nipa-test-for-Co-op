package service

import (
	"bookedticket_backend/internal/models"
	"bookedticket_backend/internal/repository"
)

type TicketService interface {
	GetTicket() ([]models.Ticket, error)
	GetTicketByID(id int) (models.Ticket, error)
	CreateTicket(ticket *models.TicketInsert) error
	UpdateTicket(id int, ut models.TicketInsert) error
}

type ticketService struct {
	repo repository.TicketRepository
}

func NewTicketService(repo repository.TicketRepository) TicketService {
	return &ticketService{repo: repo}
}

func (ts *ticketService) GetTicket() ([]models.Ticket, error) {
	return ts.repo.GetTicket()
}

func (ts *ticketService) GetTicketByID(id int) (models.Ticket, error) {
	return ts.repo.GetTicketByID(id)
}

func (ts *ticketService) CreateTicket(ticket *models.TicketInsert) error {
	return ts.repo.CreateTicket(ticket)
}

func (ts *ticketService) UpdateTicket(id int, ut models.TicketInsert) error {
	return ts.repo.UpdateTicket(id, ut)
}
