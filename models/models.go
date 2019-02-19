package models

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// Creating a structure for Taxi
type Taxi struct {
	TaxiNo                  int // Unique identification of taxi
	NumTravellingPassengers int // At present how many number of passengers in taxi
	TaxiCapacity            int // Maximum number of travellers can travel at once

}
