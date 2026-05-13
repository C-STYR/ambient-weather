package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/cstyr/weather/internal/config"
)

// This is a utility to check a single record from the API for prototyping

func main() {

	// LOAD ENV VARS
	envConfig, err := config.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	url := fmt.Sprintf("https://rt.ambientweather.net/v1/devices/%s?apiKey=%s&applicationKey=%s&limit=%d",
		envConfig.MacAddr,
		envConfig.ApiKey,
		envConfig.AppKey,
		3,
	)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error during fetch: ", err)
		return
	}

	json, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error decoding json: ", err)
		return
	}

	fmt.Println("fetched record:\n", string(json))
}
