package ambient

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type httpClient struct {
	statusCode int
	body       string
}

func (c *httpClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: c.statusCode,
		Body:       io.NopCloser(strings.NewReader(c.body)),
	}, nil
}

func TestFetchPage(t *testing.T) {

	// HAPPY PATH
	client := &Client{
		APIKey: "123",
		AppKey: "456",
		HttpClient: &httpClient{
			statusCode: 200,
			// contains 3 records
			body: `[{"dateutc":1778707500000,"tempf":51.98,"humidity":53,"windspeedmph":2.24,"windgustmph":3.36,"maxdailygust":6.93,"winddir":99,"winddir_avg10m":345,"hourlyrainin":0,"eventrainin":0,"dailyrainin":0,"weeklyrainin":0.36,"monthlyrainin":1.56,"yearlyrainin":2.6,"totalrainin":2.602,"battout":1,"tempinf":66.92,"humidityin":51,"baromrelin":29.043,"baromabsin":29.043,"feelsLike":51.98,"dewPoint":35.39,"feelsLikein":66.9,"dewPointin":48.2,"lastRain":"2026-05-11T08:00:00.000Z","date":"2026-05-13T21:25:00.000Z"},{"dateutc":1778707200000,"tempf":51.8,"humidity":55,"windspeedmph":1.34,"windgustmph":1.12,"maxdailygust":6.93,"winddir":61,"winddir_avg10m":355,"hourlyrainin":0,"eventrainin":0,"dailyrainin":0,"weeklyrainin":0.36,"monthlyrainin":1.56,"yearlyrainin":2.6,"totalrainin":2.602,"battout":1,"tempinf":67.1,"humidityin":51,"baromrelin":29.04,"baromabsin":29.04,"feelsLike":51.8,"dewPoint":36.15,"feelsLikein":67.1,"dewPointin":48.4,"lastRain":"2026-05-11T08:00:00.000Z","date":"2026-05-13T21:20:00.000Z"},{"dateutc":1778706900000,"tempf":51.44,"humidity":56,"windspeedmph":0.67,"windgustmph":2.24,"maxdailygust":6.93,"winddir":338,"winddir_avg10m":355,"hourlyrainin":0,"eventrainin":0,"dailyrainin":0,"weeklyrainin":0.36,"monthlyrainin":1.56,"yearlyrainin":2.6,"totalrainin":2.602,"battout":1,"tempinf":67.1,"humidityin":51,"baromrelin":29.046,"baromabsin":29.046,"feelsLike":51.44,"dewPoint":36.27,"feelsLikein":67.1,"dewPointin":48.4,"lastRain":"2026-05-11T08:00:00.000Z","date":"2026-05-13T21:15:00.000Z"}]`,
		},
	}

	records, err := client.FetchPage("mac", 0, 3)
	if err != nil {
		t.Fatal("failed")
	}

	if len(records) != 3 {
		t.Errorf("Expected 3 records, got %d records", len(records))
	}

	// non-200 path
	client = &Client{
		APIKey: "123",
		AppKey: "456",
		HttpClient: &httpClient{
			statusCode: 403,
			body: ``,
		},
	}

	records, err = client.FetchPage("mac", 0, 3)
	expected := "API returned status 403"
	received := err.Error()
	if expected != received {
		t.Fatal(err.Error())
	}

	if len(records) != 0 {
		t.Errorf("Expected 0 records, got %d records", len(records))
	}

	// add Get fail path
}
