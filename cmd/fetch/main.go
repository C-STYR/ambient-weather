package main

import (
	"fmt"
	"os"

	"github.com/cstyr/weather/internal/ambient"
)

func main() {
	macAddr, ok := os.LookupEnv("AMBIENT_DEVICE_MAC")
	if !ok {
		fmt.Println("mac address not set")
	}

	apiKey, ok := os.LookupEnv("AMBIENT_API_KEY")
	if !ok {
		fmt.Println("api key not set")
	}

	appKey, ok := os.LookupEnv("AMBIENT_APP_KEY")
	if !ok {
		fmt.Println("app key not set")
	}

	client := &ambient.Client{
		APIKey: apiKey,
		AppKey: appKey,
	}

	results, err := client.FetchPage(macAddr, 0, 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Results: %+v\n", results)

}
