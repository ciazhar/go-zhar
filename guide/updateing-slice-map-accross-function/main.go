package main

import "fmt"

type Trip struct {
	destination string
}

type Driver struct {
	trips []Trip
}

func (d *Driver) SetTripsBad(trips []Trip) {
	d.trips = trips
}

func (d *Driver) SetTripsGood(trips []Trip) {
	d.trips = make([]Trip, len(trips))
	copy(d.trips, trips)
}

func main() {

	fmt.Println("Bad Example")

	// Initialize the trips slice
	trips := []Trip{
		{destination: "New York"},
		{destination: "Los Angeles"},
	}

	// Initialize the Driver
	var d1 Driver
	d1.SetTripsBad(trips)

	// Print initial trips
	fmt.Println("Initial trips in d1:", d1.trips)

	// Modify the trips slice
	trips[0] = Trip{destination: "San Francisco"}

	// Print modified trips
	fmt.Println("Modified trips in original slice:", trips)
	fmt.Println("Modified trips in d1:", d1.trips)

	fmt.Println("Good Example")

	// Initialize the trips slice
	trips2 := []Trip{
		{destination: "New York"},
		{destination: "Los Angeles"},
	}

	// Initialize the Driver
	var d2 Driver
	d2.SetTripsGood(trips2)

	// Print initial trips
	fmt.Println("Initial trips in d2:", d2.trips)

	// Modify the trips slice
	trips2[0] = Trip{destination: "San Francisco"}

	// Print modified trips
	fmt.Println("Modified trips in original slice:", trips2)
	fmt.Println("Modified trips in d2:", d2.trips)

}
