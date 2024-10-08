package fr24

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
)

type AirlineRes struct {
	Version int           `json:"version"`
	Rows    []AirlineData `json:"rows"`
}

type AirlineData struct {
	Name string `json:"Name"`
	Code string `json:"Code"`
	Icao string `json:"ICAO"`
}

func GetAirlines(requester Requester) (AirlineRes, error) {
	var airlines AirlineRes

	body, err := requester(FR24_ENDPOINTS["airlines"])

	if err != nil {
		return airlines, NewFr24Error(err)
	}

	if err := json.Unmarshal(body, &airlines); err != nil {
		return airlines, NewFr24Error(err)
	}

	return airlines, nil
}

func GetAirlineLogoCdn(requester Requester, icao string, iata string) (bytes.Buffer, error) {
	/**
	* Gets a logo for a requested airline based on its icao and iata code from the CDN.
	 */
	var buf bytes.Buffer
	endpoint := fmt.Sprintf("%s/%s_%s.png", FR24_ENDPOINTS["airline_logo_cdn"], icao, iata)

	if err := createPng(&buf, endpoint, requester); err != nil {
		return buf, err
	}

	return buf, nil
}

func GetAirlineLogo(requester Requester, icao string) (bytes.Buffer, error) {
	/**
	* Get Logo from assets on Flightradar
	 */
	var buf bytes.Buffer
	endpoint := fmt.Sprintf("%s/%s_logo0.png", FR24_ENDPOINTS["airline_logo"], icao)

	if err := createPng(&buf, endpoint, requester); err != nil {
		return buf, err
	}

	return buf, nil
}

func createPng(buf *bytes.Buffer, endpoint string, requester Requester) error {
	/**
	* Sends the request to Flightradar24 for a PNG image and encodes the returned bytes into
	* a PNG logo.
	 */
	body, err := requester(endpoint)

	if err != nil {
		return NewFr24Error(err)
	}

	img, _, err := image.Decode(bytes.NewReader(body))

	if err != nil {
		return NewFr24Error(err)
	}

	// encode the bytes into a png image
	if err := png.Encode(buf, img); err != nil {
		return NewFr24Error(err)
	}

	return nil
}
