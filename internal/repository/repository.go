package repository

import "github.com/Shazil2154/-go--bookings-project/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}
