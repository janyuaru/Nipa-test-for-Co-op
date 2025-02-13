package models

import (
	"time"
)

type Ticket struct {
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Contact            string    `json:"contact"`
	Status             string    `json:"status"`
	Created_ticket_at  time.Time `json:"created_ticket_at"`
	Lastest_updated_at time.Time `json:"lastest_updated_at"`
}

type TicketInsert struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Contact     string `json:"contact"`
	Status      string `json:"status"`
}
