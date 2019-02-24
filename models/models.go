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
	Name     string
	UName    string `gorm:"type:varchar(40); not null`
	Password string `gorm:"type:varchar(40); not null`
	Email    string
	PhoneNo  string
	Wallet   int64
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
	TaxiID      uint // initialised when a taxi is booked for a given route
	Distance    int
	Duration    int
	Fare        int
	GooglePath  []GooglePath
	Source      string
	Destination string
	Status      string // can be Active or Passive
}

type Booking struct {
	DBModel
	RouteID uint
	//Route          Route
	TaxiID uint
	//Taxi           Taxi
	RiderID uint
	//Rider          Rider
	TravelDuration int
	ETA            int    // How much time a user has to wait till he gets a taxi.
	Status         string // Can be To_start, Active or Finished
}

// DashBoard Holds data to be sent when dashboard Endpoint hit
type DashBoard struct {
	Name   string
	Email  string
	Phone  string
	Wallet int64
}
