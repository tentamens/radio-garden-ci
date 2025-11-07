package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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

func PickRandonStation() string {
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
