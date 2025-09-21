package main

const configFilePath string = "config.json"
const envFilePath string = ".env"

type tickerData struct {
	Threshold float64 `json:"threshold"`
	LastPrice float64 `json:"-"`
}

type configData struct {
	APIEndpoint string                `json:"apiEndpoint"`
	LoopDelay   int                   `json:"loopDelay"`
	TGService   string                `json:"tgService"`
	BotName     string                `json:"botName"`
	BotToken    string                `json:"-"`
	ChatID      string                `json:"chatID"`
	Tickers     map[string]tickerData `json:"tickers"`
}

var appConfig configData

type rawStockData struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type stockData struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}
