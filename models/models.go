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
	Name string
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
	TaxiID		uint
	Distance 	int
	Duration    int
	Fare        int
	GooglePath  []GooglePath
	Source      string
	Destination string
}

type Booking struct {
	DBModel
	RouteID        uint
	Route          Route
	TaxiID         uint
	Taxi           Taxi
	RiderID        uint
	Rider          Rider
	ETA            int // How much time a user has to wait till he gets a taxi.
}
