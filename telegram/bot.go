package telegram

import (
	TGApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type TelegramBot struct {
	telegramBot     *TGApi.BotAPI
	telegramChatBot map[int64]*telegramChatBot
}

func (chatbot *TelegramBot) Error() string {
	//TODO implement me

	return ""
}

func NewBot(bot *TGApi.BotAPI) *TelegramBot {
	return &TelegramBot{
		telegramBot:     bot,
		telegramChatBot: make(map[int64]*telegramChatBot),
	}
}

func (tgBot *TelegramBot) StartTelegramUpdates() error {
	upd := TGApi.NewUpdate(0)
	upd.Timeout = 20
	upds := tgBot.telegramBot.GetUpdatesChan(upd)

	for update := range upds {
		if chatBot, ok := tgBot.telegramChatBot[update.FromChat().ID]; ok {
			chatBot.GetUpdateCh() <- update
		} else {
			tgBot.telegramChatBot[update.FromChat().ID] = &telegramChatBot{
				tg:                     tgBot.telegramBot,
				tradeBot:               nil,
				currChatState:          nil,
				logger:                 logrus.New(),
				newUpdatesCh:           make(chan TGApi.Update),
				closeUpdateListenerCh:  make(chan struct{}),
				closeSignalsListenerCh: make(chan struct{}),
			}
			go tgBot.telegramChatBot[update.FromChat().ID].ListenNewUpdates()
			tgBot.telegramChatBot[update.FromChat().ID].GetUpdateCh() <- update
		}
	}

	return nil
}
