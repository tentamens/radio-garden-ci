package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	radiogarden "thomasjgriffin.dev/radio-garden-cli"
)

var client *radiogarden.ClientWithResponses

func InitClient() {
	globalOpts.Server = "https://radio.garden/api"
	// Construct a client with automatic redirect following disabled to avoid
	// downloading station audio streams via the station stream command.
	c, err := radiogarden.NewClientWithResponses(
		"https://radio.garden/api")
	if err != nil {
		fmt.Println(err)
	}

	client = c
}

type APIResponse struct {
	Data struct {
		List json.RawMessage `json:"list"`
	} `json:"data"`
}

func GetLocations(ctx context.Context) {
	resp, err := client.GetAraContentPlaces(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var jsonBody APIResponse
	err = json.Unmarshal(bodyBytes, &jsonBody)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}

	fmt.Println(len(jsonBody.Data.List))

	os.WriteFile("assets/caches/locationCache.json", jsonBody.Data.List, 0o644)
}
