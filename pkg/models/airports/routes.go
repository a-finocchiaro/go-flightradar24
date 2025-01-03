package airports

import (
	"encoding/json"

	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/common"
)

type AirportRouteResponse struct {
	Arrivals   AirportRoute `json:"arrivals"`
	Departures AirportRoute `json:"departures"`
}

type AirportRoute struct {
	Country  string             `json:"id"`
	Number   AirportRouteNumber `json:"number"`
	Airports AirportRouteData   `json:"airports"`
}

// Parses all route data into an object using a custom mapping
func (a *AirportRoute) UnmarshalJSON(data []byte) error {
	// map is used here to capture the Country name and assign it
	var temp map[string]struct {
		Number   AirportRouteNumber `json:"number"`
		Airports map[string]struct {
			Name     string     `json:"name"`
			City     string     `json:"city"`
			Icao     string     `json:"Icao"`
			Position LatLongStr `json:"position"`
			Flights  map[string]struct {
				Airline AirportRouteAirline               `json:"Airline"`
				Utc     map[string]AirportRouteFlightTime `json:"utc"`
			} `json:"flights"`
		} `json:"airports"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return common.NewFr24Error(err)
	}

	// assign values into the object
	for country, routeData := range temp {
		a.Country = country
		a.Number = routeData.Number

		// get the IATA code
		for iata, airportData := range routeData.Airports {
			a.Airports.Iata = iata
			a.Airports.City = airportData.City
			a.Airports.Icao = airportData.Icao
			a.Airports.Name = airportData.Name
			a.Airports.Position = airportData.Position

			// get the flight IDs
			for id, flightData := range airportData.Flights {
				flight := AirportRouteFlightData{
					ID:      id,
					Airline: flightData.Airline,
				}

				for date, timeData := range flightData.Utc {
					timeData.Date = date
					flight.Utc = append(flight.Utc, timeData)
				}

				a.Airports.Flights = append(a.Airports.Flights, flight)
			}
		}
	}

	return nil
}

type AirportRouteNumber struct {
	Airports int `json:"airports"`
	Flights  int `json:"flights"`
}

type AirportRouteData struct {
	Iata string
	AirportRouteAirport
}

type AirportRouteAirport struct {
	Name     string                   `json:"name"`
	City     string                   `json:"city"`
	Icao     string                   `json:"icao"`
	Position LatLongStr               `json:"position"`
	Flights  []AirportRouteFlightData `json:"flights"`
}

type AirportRouteFlightData struct {
	ID      string                   `json:"id"`
	Airline AirportRouteAirline      `json:"Airline"`
	Utc     []AirportRouteFlightTime `json:"utc"`
}

type AirportRouteAirline struct {
	Name string `json:"name"`
	Iata string `json:"iata"`
	Icao string `json:"icao"`
	Url  string `json:"url"`
}

type AirportRouteFlightTime struct {
	Date string `json:"id"`
	AirportRouteFlightTimeAircraftInfo
}

type AirportRouteFlightTimeAircraftInfo struct {
	Aircraft  string `json:"aircraft"`
	Time      string `json:"time"`
	Timestamp int64  `json:"timestamp"`
	Offset    int    `json:"offset"`
}

type LatLongStr struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}
