package main

import (
	"fmt"
	"time"

	"github.com/cstyr/weather/internal/ambient"
	"github.com/cstyr/weather/internal/config"
	"github.com/cstyr/weather/internal/store"
)

func main() {
	// LOAD ALL ENV VARS
	envConfig, err := config.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &ambient.Client{
		APIKey: envConfig.ApiKey,
		AppKey: envConfig.AppKey,
	}

	store, err := store.New(envConfig.Region, envConfig.TableName)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var endDate int64

	for {
		// GET DATA FROM AMBIENT WEATHER API
		results, err := client.FetchPage(envConfig.MacAddr, endDate, 288)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		recordCount := len(results)
		if recordCount == 0 {
			break
		}
		fmt.Printf("fetched %d records\n", recordCount)

		// DB WRITE
		err = store.ProcessRecords(envConfig.StationName, 25, results)
		if err != nil {
			fmt.Println("error writing records:", err)
			return
		}

		latestResult := results[recordCount-1].DateUTC

		fmt.Printf("oldest timestamp: %d\n", latestResult)

		endDate = latestResult - 1

		// weather API is rate limited to one request per second
		time.Sleep(1 * time.Second)
	}

	fmt.Println("DONE")

}
