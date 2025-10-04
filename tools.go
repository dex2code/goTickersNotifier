package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func checkFileExists(fileName string) error {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return fmt.Errorf("%q does not exist or not accessible: %w", fileName, err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("%q is a directory not a file", fileName)
	}

	fileDescr, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("cannot open existing %q file: %w", fileName, err)
	}
	defer fileDescr.Close()

	return nil
}

func createDefaultConfig(fileName string) error {
	newConfig := configData{
		APIEndpoint: "https://api.binance.com/api/v3/ticker/price?symbols=",
		LoopDelay:   600,
		TGService:   "https://api.telegram.org/bot%s/sendMessage",
		BotName:     "🤖 cryptoBot:",
		ChatID:      "",
		Tickers: map[string]tickerData{
			"BTCUSDT": {
				Threshold: 500,
				LastPrice: 0.0,
			},
			"ETHUSDT": {
				Threshold: 100,
				LastPrice: 0.0,
			},
			"SOLUSDT": {
				Threshold: 5,
				LastPrice: 0.0,
			},
		},
	}

	jsonData, err := json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal JSON data: %w", err)
	}

	fileDescr, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cannot create default config file %q: %w", fileName, err)
	}
	defer fileDescr.Close()

	if err := os.WriteFile(fileName, jsonData, 0644); err != nil {
		return fmt.Errorf("cannot write file %q: %w", fileName, err)
	}

	return nil
}

func readConfig(fileName string, config *configData) error {
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("cannot read data from %q: %w", fileName, err)
	}

	if err := json.Unmarshal(jsonData, &config); err != nil {
		return fmt.Errorf("cannot unmarshal config data: %w", err)
	}

	return nil
}

func createDefaultEnv(fileName string) error {
	envContent := "BOT_TOKEN=\"\""
	if err := os.WriteFile(fileName, []byte(envContent), 0600); err != nil {
		return fmt.Errorf("cannot create %q file: %w", fileName, err)
	}

	return nil
}
