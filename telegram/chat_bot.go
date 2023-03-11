package telegram

import (
	"fmt"
	TGApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"someshit/cmd/trade"
)

type ChatState struct {
	isCallback bool
	isCommand  bool
	value      string // хранит команду или выбранный вариант для callback
	figi       string
}

type telegramChatBot struct {
	tg                     *TGApi.BotAPI
	tradeBot               *trade.TradeBot
	currChatState          *ChatState
	logger                 *logrus.Logger
	newUpdatesCh           chan TGApi.Update
	closeUpdateListenerCh  chan struct{}
	closeSignalsListenerCh chan struct{}
}

func (chatbot *telegramChatBot) ListenNewUpdates() {
	for {
		select {
		case upd := <-chatbot.newUpdatesCh:
			err := chatbot.UpdateHandler(upd)
			if err != nil {
				chatbot.logger.Error(err)
			}
		case <-chatbot.closeUpdateListenerCh:
			return
		}
	}
}

func (chatbot *telegramChatBot) ListenSignalsFromTinkoffBot(id int64) {
	//defer wg.Done()
	for {
		select {
		case signal := <-chatbot.tradeBot.GetWorkersChangesCh():
			if signal.SignalsType == trade.Err_type {
				chatbot.tradeBot.StopWorker(signal.WorkerId)
			}
			msg := TGApi.NewMessage(id, signal.Description)
			_, err := chatbot.tg.Send(msg)
			if err != nil {
				chatbot.logger.Error(err)
			}
			if signal.Img != nil {
				for _, k := range signal.Img {
					photo := TGApi.FileBytes{
						Name:  fmt.Sprintf("График от %d", signal.WorkerId),
						Bytes: k,
					}
					ph := TGApi.NewPhoto(id, photo)
					_, err = chatbot.tg.Send(ph)
					if err != nil {
						chatbot.logger.Error(err)
					}
				}

			}

		case <-chatbot.closeSignalsListenerCh:
			return
		}
	}
}

func (chatbot *telegramChatBot) StopTinkoffBot(id int64) {
	chatbot.tradeBot.StopBot()                   //останавливаем воркеров
	chatbot.closeSignalsListenerCh <- struct{}{} //останавливаем горутину, которая слушает сигналы воркеров от этого бота
	chatbot.tradeBot = nil
	close(chatbot.closeSignalsListenerCh)
	return
}

func (chatbot *telegramChatBot) StartBot(id int64) {
	//defer chatbot.botWg.Done()
	chatbot.closeSignalsListenerCh = make(chan struct{})
	//wg := new(sync.WaitGroup)
	//wg.Add(2)
	go chatbot.ListenSignalsFromTinkoffBot(id)
	go chatbot.tradeBot.StartNewWorkers()
	//go chatbot.ListenNewUpdates()
	//wg.Wait()

}

func (chatbot *telegramChatBot) UpdateHandler(update TGApi.Update) error {
	//если обновление - сообщение
	if update.Message != nil {
		//если сообщение - команда
		if update.Message.IsCommand() {
			//запоминаем текущую команду в чате
			chatbot.currChatState = &ChatState{
				isCallback: false,
				isCommand:  true,
				value:      update.Message.Command(),
			}

			//запускаем обработчик команды
			if err := chatbot.handleCommands(update.Message); err != nil {
				//chatbot.logger.Error(err)
				return err
			}
		} else {
			if chatbot.currChatState == nil {
				//TODO
				msg := TGApi.NewMessage(update.Message.Chat.ID, "Сначала введите команду")
				_, err := chatbot.tg.Send(msg)
				if err != nil {
					//chatbot.logger.Error(err)
					return err
				}
				return nil
			}
			//если получено сообщение и до этого в чате была введена команда
			if chatbot.currChatState.isCommand {
				err := chatbot.handleTextMessageAfterCommand(update.Message)
				if err != nil {
					//chatbot.logger.Error(err)
					return err
				}

				//если было получено сообщение и до этого в чате был выбран вариант в callback
			} else if chatbot.currChatState.isCallback {
				err := chatbot.handleTextMessageAfterCallback(update.Message)
				if err != nil {
					//chatbot.logger.Error(err)
					return err
				}

			}

		}
		//если обновление - callback
	} else if update.CallbackQuery != nil && chatbot.currChatState != nil {
		err := chatbot.handleCallbackUpdate(update.FromChat().ID, update.CallbackQuery)
		if err != nil {
			//chatbot.logger.Error(err)
			return err
		}
	}
	return nil
}

func (chatbot *telegramChatBot) GetUpdateCh() chan TGApi.Update {
	return chatbot.newUpdatesCh
}
