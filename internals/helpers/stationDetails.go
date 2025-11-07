package helpers

import (
	"encoding/json"
	"net/http"
)

type Country struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Place struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Data struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Website string  `json:"website"`
	Secure  bool    `json:"secure"`
	Place   Place   `json:"place"`
	Country Country `json:"country"`
}

type Response struct {
	Data Data `json:"data"`
}

type StationDetails struct {
	Title   string
	Website string
	Country string
	City    string
}

func GetStationDetails(stationID string) (StationDetails, error) {
	url := "https://radio.garden/api/ara/content/channel/" + stationID

	details := StationDetails{}

	resp, err := http.Get(url)
	if err != nil {
		return details, err
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return details, err
	}

	details = StationDetails{
		Title:   response.Data.Title,
		Website: response.Data.Website,
		Country: response.Data.Country.Title,
		City:    response.Data.Place.Title,
	}

	return details, nil
}
