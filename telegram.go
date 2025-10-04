package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func composeTgMessage(ticker string, lastPrice float64, curPrice float64) string {
	ind := ""
	if curPrice < lastPrice {
		ind = "🔴"
	} else {
		ind = "🟢"
	}

	p := message.NewPrinter(language.English)

	if lastPrice == 0.0 {
		return p.Sprintf("%s %s %s: $%g", appConfig.BotName, ind, ticker, curPrice)
	}

	return p.Sprintf("%s %s %s: $%g -> $%g", appConfig.BotName, ind, ticker, lastPrice, curPrice)
}

func sendTgMessage(text string) error {
	tgAPI := fmt.Sprintf(appConfig.TGService, appConfig.BotToken)

	formData := url.Values{}
	formData.Add("chat_id", appConfig.ChatID)
	formData.Add("text", text)

	resp, err := http.PostForm(tgAPI, formData)
	if err != nil {
		return fmt.Errorf("cannot send TG message: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read TG service answer: %w", err)
	}

	var tgResponse TelegramResponse
	if err := json.Unmarshal(body, &tgResponse); err != nil {
		return fmt.Errorf("cannot unmarshal TG answer: %w", err)
	}

	if !tgResponse.OK {
		return fmt.Errorf("telegram API error: %w", err)
	}

	return nil
}
