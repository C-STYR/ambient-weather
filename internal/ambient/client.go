package ambient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type Client struct {
	APIKey     string
	AppKey     string
	HttpClient HTTPClient
}

type Record struct {
	DateUTC      int64   `json:"dateutc"       dynamodbav:"dateutc"`
	Date         string  `json:"date"       dynamodbav:"date"`
	TempF        float64 `json:"tempf"         dynamodbav:"tempf"`
	Humidity     float64 `json:"humidity"      dynamodbav:"humidity"`
	WindSpeedMph float64 `json:"windspeedmph"  dynamodbav:"windspeedmph"`
	WindGustMph  float64 `json:"windgustmph"   dynamodbav:"windgustmph"`
	MaxDailyGust float64 `json:"maxdailygust" dynamodbav:"maxdailygust"`
	WindDir      float64 `json:"winddir"   dynamodbav:"winddir"`
	DailyRainIn  float64 `json:"dailyrainin"   dynamodbav:"dailyrainin"`
	LastRain     string  `json:"lastRain"   dynamodbav:"lastRain"`
	BaromRelIn   float64 `json:"baromrelin"    dynamodbav:"baromrelin"`
	FeelsLike    float64 `json:"feelsLike"     dynamodbav:"feelsLike"`
	DewPoint     float64 `json:"dewPoint"      dynamodbav:"dewPoint"`
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

	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// handle error from Ambient Weather API
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// decode the json response body into []Record
	var records []Record
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, err
	}

	return records, nil
}
