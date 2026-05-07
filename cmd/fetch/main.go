package main

import (
	"fmt"
	"os"

	"github.com/cstyr/weather/internal/ambient"
	"github.com/cstyr/weather/internal/store"
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

	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		fmt.Println("region not set")
	}

	tableName, ok := os.LookupEnv("WEATHER_TABLE_DDB")
	if !ok {
		fmt.Println("tablename not set")
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

	store, err := store.New(region, tableName)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	err = store.WriteRecords("112thFabShop_1", results)
	if err != nil {
		fmt.Println("error writing records:", err)
		return
	}
	fmt.Println("DONE")

}
