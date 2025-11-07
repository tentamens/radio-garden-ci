package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type Location struct {
	Size    int       `json:"size"`
	ID      string    `json:"id"`
	Geo     []float64 `json:"geo"`
	URL     string    `json:"url"`
	Boost   bool      `json:"boost"`
	Title   string    `json:"title"`
	Country string    `json:"country"`
}

func PickRandonPlace() string {
	const filePath = "assets/caches/locationCache.json"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var jsonData []Location
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	randmonNumber := rand.Intn(len(jsonData))

	pickedItem := jsonData[randmonNumber]

	splitURL := strings.Split(pickedItem.URL, "/")

	urlCode := splitURL[len(splitURL)-1]

	return urlCode
}

func PickStation(ctx context.Context, stationID string) (string, error) {
	resp, err := http.Get("https://radio.garden/api/ara/content/page/" + stationID)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	// 3. Read the entire body content into a byte slice
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	var result map[string]any
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal root object: %w", err)
	}

	// 2. Access 'content' (array) and assert its type
	data, ok := result["data"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("data not found or is empty")
	}

	content, ok := data["content"].([]any)
	if !ok || len(content) == 0 {
		return "", fmt.Errorf("content not found or is empty")
	}
	// 3. Access first item of 'content' (object) and assert its type
	firstSection, ok := content[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("first section is not an object")
	}

	// 4. Access 'items' (array) and assert its type
	items, ok := firstSection["items"].([]any)
	if !ok || len(items) == 0 {
		return "", fmt.Errorf("items array not found or is empty")
	}

	// 5. Access first item of 'items' (object) and assert its type
	firstItem, ok := items[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("first item is not an object")
	}

	// 6. Access 'page' (object) and assert its type
	page, ok := firstItem["page"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("page object not found")
	}

	// 7. Access 'stream' (string) and assert its type
	streamURL, ok := page["url"].(string)
	if !ok {
		return "", fmt.Errorf("stream URL not found or is not a string")
	}

	splitURL := strings.Split(streamURL, "/")

	urlCode := splitURL[len(splitURL)-1]

	return urlCode, nil
}
