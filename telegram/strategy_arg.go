package telegram

import (
	"final/internal/strategy"
	"final/internal/trade"
	"fmt"
	TGApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"reflect"
	"strconv"
)

func (chatbot *telegramChatBot) doubleEMA_Arg(msg *TGApi.Message) error {
	strat := new(strategy.DoubleEMA)
	params := StrategyParamParser(msg.Text)
	if len(params) == 0 {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Нет аргументов для индикатора или неверное форматирование"))
		_, err := chatbot.tg.Send(RespMsg)
		return err
	}
	val := reflect.TypeOf(strat).Elem()
	numbers := 0
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Tag.Get("reflect") != "-" {
			numbers++
		}
	}
	if len(params) != numbers+1 {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Неверное количество аргументов для "+
			"индикатора, ожидалось %d, получено %d", numbers+1, len(params)))
		_, err := chatbot.tg.Send(RespMsg)
		return err
	}

	var err error
	errResp := TGApi.NewMessage(msg.Chat.ID, "Неверный формат данных")
	figi := params[0]
	strat.ShortWindow, err = strconv.Atoi(params[1])
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}
	strat.LongWindow, err = strconv.Atoi(params[2])
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}
	strat.CandleIntervalHours, err = strconv.Atoi(params[3])
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}
	/*strat.AnalyzeIntervalHours, err = strconv.Atoi(params[4])
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}*/
	strat.WhichEMAHigher = false

	worker := trade.NewWorkerConfig(figi, strat)
	ch := chatbot.tradeBot.GetNewWorkersConfigCh()
	ch <- worker
	return err
}

func (chatbot *telegramChatBot) bb_Arg(msg *TGApi.Message) error {
	strat := new(strategy.BollingerBands)
	params := StrategyParamParser(msg.Text)
	if len(params) == 0 {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Нет аргументов для индикатора или неверное форматирование"))
		_, err := chatbot.tg.Send(RespMsg)
		return err
	}
	val := reflect.TypeOf(strat).Elem()
	numbers := 0
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Tag.Get("reflect") != "-" {
			numbers++
		}
	}
	if len(params) != numbers+1 {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Неверное количество аргументов для "+
			"индикатора, ожидалось %d, получено %d", numbers+1, len(params)))
		_, err := chatbot.tg.Send(RespMsg)
		return err
	}

	var err error
	errResp := TGApi.NewMessage(msg.Chat.ID, "Неверный формат данных")
	figi := params[0]
	strat.Window, err = strconv.Atoi(params[1])
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}
	strat.Sigma, err = strconv.ParseFloat(params[2], 64)
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}
	strat.CandleIntervalHours, err = strconv.Atoi(params[3])
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}

	worker := trade.NewWorkerConfig(figi, strat)
	ch := chatbot.tradeBot.GetNewWorkersConfigCh()
	ch <- worker
	return err
}

func (chatbot *telegramChatBot) macd_Arg(msg *TGApi.Message) error {
	strat := new(strategy.MACD)
	params := StrategyParamParser(msg.Text)
	if len(params) == 0 {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Нет аргументов для индикатора или неверное форматирование"))
		_, err := chatbot.tg.Send(RespMsg)
		return err
	}
	val := reflect.TypeOf(strat).Elem()
	numbers := 0
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Tag.Get("reflect") != "-" {
			numbers++
		}
	}
	if len(params) != numbers+1 {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Неверное количество аргументов для "+
			"индикатора, ожидалось %d, получено %d", numbers+1, len(params)))
		_, err := chatbot.tg.Send(RespMsg)
		return err
	}

	var err error
	errResp := TGApi.NewMessage(msg.Chat.ID, "Неверный формат данных")
	figi := params[0]

	strat.ShortWindow, err = strconv.Atoi(params[1])
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}
	strat.LongWindow, err = strconv.Atoi(params[2])
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}
	strat.CandleIntervalHours, err = strconv.Atoi(params[3])
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}

	worker := trade.NewWorkerConfig(figi, strat)
	ch := chatbot.tradeBot.GetNewWorkersConfigCh()
	ch <- worker
	return err
}

func (chatbot *telegramChatBot) rsi_Arg(msg *TGApi.Message) error {
	strat := new(strategy.RSI)
	params := StrategyParamParser(msg.Text)
	if len(params) == 0 {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Нет аргументов для индикатора или неверное форматирование"))
		_, err := chatbot.tg.Send(RespMsg)
		return err
	}
	val := reflect.TypeOf(strat).Elem()
	numbers := 0
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Tag.Get("reflect") != "-" {
			numbers++
		}
	}
	if len(params) != numbers+1 {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Неверное количество аргументов для "+
			"индикатора, ожидалось %d, получено %d", numbers+1, len(params)))
		_, err := chatbot.tg.Send(RespMsg)
		return err
	}

	var err error
	errResp := TGApi.NewMessage(msg.Chat.ID, "Неверный формат данных")
	figi := params[0]

	strat.HighBorder, err = strconv.ParseFloat(params[1], 64)
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}
	strat.LowBorder, err = strconv.ParseFloat(params[2], 64)
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}
	strat.CandleIntervalHours, err = strconv.Atoi(params[4])
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}
	strat.Window, err = strconv.Atoi(params[3])
	if err != nil {
		_, err = chatbot.tg.Send(errResp)
		return err
	}

	worker := trade.NewWorkerConfig(figi, strat)
	ch := chatbot.tradeBot.GetNewWorkersConfigCh()
	ch <- worker
	return err
}
