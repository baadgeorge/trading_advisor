package telegram

import (
	"context"
	"fmt"
	TGApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

// структура телеграм бота
type TelegramBot struct {
	telegramBot     *TGApi.BotAPI
	telegramChatBot map[int64]*telegramChatBot
	stopChatBotCh   chan int64
}

// функция создания экземпляра телеграм бота
func NewBot(bot *TGApi.BotAPI) *TelegramBot {
	return &TelegramBot{
		telegramBot:     bot,
		telegramChatBot: make(map[int64]*telegramChatBot),
		stopChatBotCh:   make(chan int64),
	}
}

// метод получения обновления от телеграма для всех пользователей
func (tgBot *TelegramBot) StartTelegramUpdates() error {
	upd := TGApi.NewUpdate(0)
	upd.Timeout = 20
	upds := tgBot.telegramBot.GetUpdatesChan(upd)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go tgBot.ListenStopChatBot(ctx)
	for update := range upds {
		if chatBot, ok := tgBot.telegramChatBot[update.FromChat().ID]; ok {
			chatBot.GetUpdateCh() <- update
			fmt.Printf("use old chatbot %v\n", tgBot.telegramChatBot[update.FromChat().ID])
		} else {
			log := logrus.New()
			log.SetReportCaller(true)
			tgBot.telegramChatBot[update.FromChat().ID] = &telegramChatBot{
				id:                     update.FromChat().ID,
				tg:                     tgBot.telegramBot,
				tradeBot:               nil,
				currChatState:          nil,
				logger:                 log,
				newUpdatesCh:           make(chan TGApi.Update),
				closeUpdateListenerCh:  make(chan struct{}),
				closeSignalsListenerCh: make(chan struct{}),
				stopChatBotCh:          tgBot.stopChatBotCh,
			}
			fmt.Printf("create new chatbot %v\n", tgBot.telegramChatBot[update.FromChat().ID])
			go tgBot.telegramChatBot[update.FromChat().ID].ListenNewUpdates()
			tgBot.telegramChatBot[update.FromChat().ID].GetUpdateCh() <- update
		}
	}
	return nil
}

func (tgBot *TelegramBot) ListenStopChatBot(ctx context.Context) {
	for {
		select {
		case id := <-tgBot.stopChatBotCh:
			delete(tgBot.telegramChatBot, id)
			println(" listen stop signal")
		case <-ctx.Done():
			return
		}
	}
}
