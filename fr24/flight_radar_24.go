/*
Queries Flightradar24 for the most tracked aircraft.
*/
package fr24

import (
	"encoding/json"
	"log"
)

type Fr24MostTrackedRes struct {
	Version     string                `json:"version"`
	Update_time float64               `json:"update_time"`
	Data        []Fr24MostTrackedData `json:"data"`
}

type Fr24MostTrackedData struct {
	Flight_id     string `json:"flight_id"`
	Flight        string `json:"flight"`
	Callsign      string `json:"callsign"`
	Squawk        string `json:"squawk"`
	Clicks        int    `json:"clicks"`
	From_iata     string `json:"from_iata"`
	From_city     string `json:"from_city"`
	To_iata       string `json:"to_iata"`
	To_city       string `json:"to_city"`
	Model         string `json:"model"`
	Aircraft_type string `json:"type"`
}

type (
	Requester func(string) ([]byte, error)
)

func GetFR24MostTracked(requester Requester) (Fr24MostTrackedRes, error) {
	var most_tracked Fr24MostTrackedRes
	body, err := requester(FR24_ENDPOINTS["most_tracked"])

	if err != nil {
		log.Fatalln(err)
		return most_tracked, err
	}

	if err := json.Unmarshal(body, &most_tracked); err != nil {
		return most_tracked, NewFr24Error(err)
	}

	return most_tracked, nil
}
