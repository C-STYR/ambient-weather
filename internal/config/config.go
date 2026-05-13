package config

import (
	"fmt"
	"os"
)

type Config struct {
	MacAddr     string
	ApiKey      string
	AppKey      string
	Region      string
	TableName   string
	StationName string
}

func Load() (*Config, error) {

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

	return &Config{
		MacAddr:     macAddr,
		ApiKey:      apiKey,
		AppKey:      appKey,
		Region:      region,
		TableName:   tableName,
		StationName: stationName,
	}, nil
}
