package models

import (
	"time"
)

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type DBModel struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type Rider struct {
	DBModel
	Name        string
	Source      string
	Destination string
	//InTaxiNum   int // In which taxi user is riding in
	MaximumWaitingTime int // Time in minutes, later we can try time
}

type Taxi struct {
	DBModel
	NumberPlate            string
	CarType                string
	CurrentLocation        string
	Capacity               int
	CurrentNumOfTravellers int
	Status                 string // Can be one of Free, Active, Full
}

type GooglePath struct {
	DBModel
	RouteID uint
	Path    string
}

type Route struct {
	DBModel
	Source      string
	Destination string
	GooglePath  []GooglePath
}

type Booking struct {
	DBModel
	RouteID        uint
	Route          Route
	TaxiID         uint
	Taxi           Taxi
	RiderID        uint
	Rider          Rider
	Price          int
	TravelDuration int // To say how much time it takes to complete the trip
	ETA            int // How much time a user has to wait till he gets a taxi.
}
