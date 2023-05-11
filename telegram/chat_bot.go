package telegram

import (
	"final/internal/trade"
	"fmt"
	TGApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

// структура состояния чата
type ChatState struct {
	isCallback bool
	isCommand  bool
	value      string // хранит команду или выбранный вариант для callback
	//figi       string
}

// структура чат бота для конкретного пользователя
type telegramChatBot struct {
	id                     int64
	tg                     *TGApi.BotAPI
	tradeBot               *trade.TradeBot
	currChatState          *ChatState
	logger                 *logrus.Logger
	newUpdatesCh           chan TGApi.Update
	closeUpdateListenerCh  chan struct{}
	closeSignalsListenerCh chan struct{}
	stopChatBotCh          chan int64
}

// метод получения обновлений от телеграма в конкретном чате
func (chatbot *telegramChatBot) ListenNewUpdates() {
	for {
		select {
		case upd := <-chatbot.newUpdatesCh:
			err := chatbot.UpdateHandler(upd)
			if err != nil {
				chatbot.logger.Error(err)
			}
		case <-chatbot.closeUpdateListenerCh:
			println("got stopChatbot")
			return
		}
	}
}

// метод получения сигналов от торгового бота
func (chatbot *telegramChatBot) ListenSignalsFromTradeBot(id int64) {
	//defer wg.Done()
	for {
		select {
		case <-chatbot.closeSignalsListenerCh:
			return
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

		}
	}
}

// метод остановки бота
func (chatbot *telegramChatBot) StopTradeBot() {
	chatbot.tradeBot.StopBot()                   //останавливаем воркеров
	chatbot.closeSignalsListenerCh <- struct{}{} //останавливаем горутину, которая слушает сигналы воркеров от этого бота
	chatbot.tradeBot = nil
	close(chatbot.closeSignalsListenerCh)
	println("1")
	return
}

func (chatbot *telegramChatBot) StopChatBot() {
	chatbot.closeUpdateListenerCh <- struct{}{}
	println("2")
	close(chatbot.closeUpdateListenerCh)
	println("3")
	chatbot.stopChatBotCh <- chatbot.id
	println("4")
	return
}

// метод запуска бота
func (chatbot *telegramChatBot) StartBot(id int64) {
	//defer chatbot.botWg.Done()
	//chatbot.closeSignalsListenerCh = make(chan struct{})
	//wg := new(sync.WaitGroup)
	//wg.Add(2)
	go chatbot.ListenSignalsFromTradeBot(id)
	go chatbot.tradeBot.TradeBotRun()
	//go chatbot.ListenNewUpdates()
	//wg.Wait()
	return
}

// метод получения обновлений от телеграм
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

		err := chatbot.handleCallback(update.FromChat().ID, update.CallbackQuery)
		if err != nil {
			//chatbot.logger.Error(err)
			return err
		}
		chatbot.currChatState = &ChatState{
			isCallback: true,
			isCommand:  false,
			value:      update.CallbackQuery.Data,
		}
	}
	return nil
}

// метод возвращающий канал получения обновлений
func (chatbot *telegramChatBot) GetUpdateCh() chan TGApi.Update {
	return chatbot.newUpdatesCh
}

func (chatbot *TelegramBot) Error() string {
	//TODO implement me

	return ""
}
