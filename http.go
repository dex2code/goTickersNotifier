package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func getStockData() ([]stockData, error) {
	resp, err := http.Get(appConfig.APIEndpoint)
	if err != nil {
		return nil, fmt.Errorf("cannot get API answer: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read API answer: %w", err)
	}

	var rawStockAnswer []rawStockData = nil

	err = json.Unmarshal(body, &rawStockAnswer)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal API answer: %w", err)
	}

	var localStockAnswer []stockData = nil

	for _, value := range rawStockAnswer {
		valuePrice, err := strconv.ParseFloat(value.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert API price: %w", err)
		}
		localStockAnswer = append(localStockAnswer, stockData{
			Symbol: value.Symbol,
			Price:  valuePrice,
		})
	}

	return localStockAnswer, nil
}
