package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Tickers Notifier starting...")
	if err := checkFileExists(envFilePath); err != nil {
		log.Printf("Env file error: %v", err)
		log.Printf("Trying to create %q file...", envFilePath)
		if err := createDefaultEnv(envFilePath); err != nil {
			log.Fatalf("Env file error: %v", err)
		} else {
			log.Printf("Default Env file created: %q", envFilePath)
		}
	} else {
		log.Printf("Found readable env file: %q", envFilePath)
	}

	if err := checkFileExists(configFilePath); err != nil {
		log.Printf("Config file error: %v", err)
		log.Printf("Trying to create %q file...", configFilePath)
		if err := createDefaultConfig(configFilePath); err != nil {
			log.Fatalf("Config file error: %v", err)
		} else {
			log.Printf("Default config file created: %q", configFilePath)
		}
	} else {
		log.Printf("Found readable config file: %q", configFilePath)
	}

	if err := godotenv.Load(envFilePath); err != nil {
		log.Fatalf("Cannot load %q file: %v", envFilePath, err)
	}
	appConfig.BotToken = os.Getenv("BOT_TOKEN")
	if appConfig.BotToken == "" {
		log.Fatalf("Please edit %q file to set secret bot token", envFilePath)
	}

	if err := readConfig(configFilePath, &appConfig); err != nil {
		log.Fatalf("No valid config: %v", err)
	} else {
		log.Printf("Readed config file successfully!")
	}

	if appConfig.ChatID == "" {
		log.Fatalf("Please edit %q file to set the correct Chat ID", configFilePath)
	}
	log.Printf("ChatID: %q", appConfig.ChatID)

	requestSuffix := []string{}
	log.Printf("Tickers:")
	for tickerName, tickerData := range appConfig.Tickers {
		requestSuffix = append(requestSuffix, "\""+tickerName+"\"")
		log.Printf("\t- %q with treshold %f", tickerName, tickerData.Threshold)
	}
	appConfig.APIEndpoint += "[" + strings.Join(requestSuffix, ",") + "]"
	log.Printf("API Endpoint: %s", appConfig.APIEndpoint)

	if err := sendTgMessage(fmt.Sprintf("%s launched to watch the following tickers:\n%s", appConfig.BotName, strings.Join(requestSuffix, ", "))); err != nil {
		log.Printf("Error sending TG message: %v", err)
	}

	log.Printf("Entering main loop with delay %d seconds...", appConfig.LoopDelay)
	firstShot := true
	for {
		if firstShot {
			firstShot = false
		} else {
			time.Sleep(time.Duration(appConfig.LoopDelay) * time.Second)
		}

		stockAnswer, err := getStockData()
		if err != nil {
			log.Printf("Cannot operate API data: %v", err)
			continue
		}

		log.Printf("Got correct answer from API: %v", stockAnswer)

		for _, stockValue := range stockAnswer {
			ticker, exists := appConfig.Tickers[stockValue.Symbol]
			if exists {
				if math.Abs(stockValue.Price-ticker.LastPrice) > ticker.Threshold {
					log.Printf("%q price change (%f -> %f) has exceeded the threshold value %f", stockValue.Symbol, ticker.LastPrice, stockValue.Price, ticker.Threshold)
					if err := sendTgMessage(composeTgMessage(stockValue.Symbol, ticker.LastPrice, stockValue.Price)); err != nil {
						log.Printf("Error sending TG message: %v", err)
						continue
					} else {
						ticker.LastPrice = stockValue.Price
						appConfig.Tickers[stockValue.Symbol] = ticker
						log.Printf("TG message sent successfully. New %q price: %f", stockValue.Symbol, ticker.LastPrice)
					}
				} else {
					log.Printf("%q price change (%f -> %f) within the treshold value %f", stockValue.Symbol, ticker.LastPrice, stockValue.Price, ticker.Threshold)
				}
			} else {
				log.Printf("API returned an unexpected ticker name %q - not found in local data", stockValue.Symbol)
			}
		}
		log.Printf("Go sleep for %d seconds...", appConfig.LoopDelay)
	}

}
