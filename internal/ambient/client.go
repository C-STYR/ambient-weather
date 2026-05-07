package ambient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	APIKey string
	AppKey string
}

type Record struct {
	DateUTC      int64   `json:"dateutc"`
	TempF        float64 `json:"tempf"`
	Humidity     float64 `json:"humidity"`
	WindSpeedMph float64 `json:"windspeedmph"`
	WindGustMph  float64 `json:"windgustmph"`
	DailyRainIn  float64 `json:"dailyrainin"`
	BaromRelIn   float64 `json:"baromrelin"`
	FeelsLike    float64 `json:"feelsLike"`
	DewPoint     float64 `json:"dewPoint"`
}

func (c *Client) FetchPage(macAddress string, endDate int64, limit int) ([]Record, error) {
	// build the url
	url := fmt.Sprintf("https://rt.ambientweather.net/v1/devices/%s?apiKey=%s&applicationKey=%s&limit=%d",
		macAddress,
		c.APIKey,
		c.AppKey,
		limit,
	)

	// add the endDate to query string if nonzero
	if endDate != 0 {
		url += fmt.Sprintf("&endDate=%d", endDate)
	}

	// make a get request using net/http
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// handle error from API
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// decode the json response body into []Record
	var records []Record
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, err
	}

	// return the slice and error
	return records, nil
}
