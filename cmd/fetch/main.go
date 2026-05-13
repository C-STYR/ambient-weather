package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cstyr/weather/internal/ambient"
	"github.com/cstyr/weather/internal/store"
)

type EnvConfig struct {
	macAddr     string
	apiKey      string
	appKey      string
	region      string
	tableName   string
	stationName string
}

func main() {
	envConfig, err := loadEnvConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &ambient.Client{
		APIKey: envConfig.apiKey,
		AppKey: envConfig.appKey,
	}

	store, err := store.New(envConfig.region, envConfig.tableName)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var endDate int64

	for {
		results, err := client.FetchPage(envConfig.macAddr, endDate, 288)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		// fmt.Printf("Results: %+v\n", results)
		recordCount := len(results)
		if recordCount == 0 {
			break
		}
		fmt.Printf("fetched %d records", recordCount)

		// for single writes only
		// err = store.WriteRecords(envConfig.stationName, results)
		
		err = store.ProcessRecords(envConfig.stationName, 25, results)
		if err != nil {
			fmt.Println("error writing records:", err)
			return
		}

		latestResult := results[recordCount-1].DateUTC

		fmt.Println("oldest timestamp: ", latestResult)

		endDate = latestResult - 1

		// weather API is rate limited to one request per second
		time.Sleep(1 * time.Second)
	}

	fmt.Println("DONE")

}

func loadEnvConfig() (*EnvConfig, error) {

	macAddr, ok := os.LookupEnv("AMBIENT_DEVICE_MAC")
	if !ok {
		err := fmt.Errorf("mac address not set")
		return nil, err
	}

	apiKey, ok := os.LookupEnv("AMBIENT_API_KEY")
	if !ok {
		err := fmt.Errorf("api key not set")
		return nil, err
	}

	appKey, ok := os.LookupEnv("AMBIENT_APP_KEY")
	if !ok {
		err := fmt.Errorf("app key not set")
		return nil, err
	}

	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		err := fmt.Errorf("region not set")
		return nil, err
	}

	tableName, ok := os.LookupEnv("WEATHER_TABLE_DDB")
	if !ok {
		err := fmt.Errorf("tablename not set")
		return nil, err
	}

	stationName, ok := os.LookupEnv("STATION_NAME")
	if !ok {
		err := fmt.Errorf("station not set")
		return nil, err
	}

	return &EnvConfig{
		macAddr:     macAddr,
		apiKey:      apiKey,
		appKey:      appKey,
		region:      region,
		tableName:   tableName,
		stationName: stationName,
	}, nil
}
