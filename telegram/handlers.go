package telegram

import (
	"errors"
	"fmt"
	TGApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"someshit/cmd/trade"
	"someshit/internal/strategy"
	"strings"
)

// HANDLERS
func (chatbot *telegramChatBot) handleTextMessageAfterCommand(msg *TGApi.Message) error {
	switch chatbot.currChatState.value {
	case "tinkoff_token":
		err := chatbot.tinkoff_token_Arg(msg)
		return err
	case "stop_worker":
		err := chatbot.stop_worker_Arg(msg)
		return err
	default:
		resp := TGApi.NewMessage(msg.Chat.ID, "У этой команды нет аргументов")
		_, err := chatbot.tg.Send(resp)
		return err
	}
}

func (chatbot *telegramChatBot) handleTextMessageAfterCallback(msg *TGApi.Message) error {
	if chatbot.tradeBot == nil {
		resp := TGApi.NewMessage(msg.Chat.ID, "Сначала нужно ввести токен! Воспользуйтесь командой /tinkoff_token")
		_, err := chatbot.tg.Send(resp)
		return err
	}
	switch chatbot.currChatState.value {
	case "doubleEMA":
		return chatbot.doubleEMA_Arg(msg)
	case "bollingerBands":
		return chatbot.bb_Arg(msg)
	case "bonds", "etfs", "shares", "futures", "currencies":
		err := chatbot.get_figi_Arg(msg)
		return err
	default:
		return nil
	}
}

func (chatbot *telegramChatBot) handleCommands(msg *TGApi.Message) error {
	switch msg.Command() {
	case "start":
		return chatbot.start_Command(msg)
	case "tinkoff_token":
		return chatbot.tinkoff_token_Command(msg)
	case "worker_list":
		if chatbot.tradeBot == nil {
			resp := TGApi.NewMessage(msg.Chat.ID, "Сначала нужно ввести токен! Воспользуйтесь командой /tinkoff_token")
			_, err := chatbot.tg.Send(resp)
			return err
		}
		return chatbot.worker_list_Command(msg)
	case "new_worker":
		//только комнада /tinkoff_token_Arg добавляет пользователя(бота)
		if chatbot.tradeBot == nil {
			resp := TGApi.NewMessage(msg.Chat.ID, "Сначала нужно ввести токен! Воспользуйтесь командой /tinkoff_token")
			_, err := chatbot.tg.Send(resp)
			return err
		}
		return chatbot.new_worker_Command(msg)
	case "stop_worker":
		if chatbot.tradeBot == nil {
			resp := TGApi.NewMessage(msg.Chat.ID, "Сначала нужно ввести токен! Воспользуйтесь командой /tinkoff_token")
			_, err := chatbot.tg.Send(resp)
			return err
		}
		return chatbot.stop_worker_Command(msg)
	case "delete_bot":
		return chatbot.delete_bot_Command(msg)
	case "help":
		return nil
	case "get_figi":
		if chatbot.tradeBot == nil {
			resp := TGApi.NewMessage(msg.Chat.ID, "Сначала нужно ввести токен! Воспользуйтесь командой /tinkoff_token")
			_, err := chatbot.tg.Send(resp)
			return err
		}
		return chatbot.get_figi_Command(msg)
	default:
		resp := TGApi.NewMessage(msg.Chat.ID, "Неизвестная команда")
		_, err := chatbot.tg.Send(resp)
		return err
	}
}

func (chatbot *telegramChatBot) handleCallbackUpdate(chatId int64, callbackQuery *TGApi.CallbackQuery) error {
	callback := TGApi.NewCallback(callbackQuery.ID, callbackQuery.Data)
	if _, err := chatbot.tg.Request(callback); err != nil {
		chatbot.logger.Warn(err)
	}
	switch chatbot.currChatState.value {
	case "get_figi":
		if chatbot.tradeBot == nil {
			resp := TGApi.NewMessage(chatId, "Сначала нужно ввести токен! Воспользуйтесь командой /tinkoff_token")
			_, err := chatbot.tg.Send(resp)
			return err
		}
		entry := chatbot.currChatState
		entry.isCallback = true
		entry.isCommand = false
		entry.value = callbackQuery.Data

		fmt.Print(callbackQuery.Data)
		//нужна функция которая возьмет выбор из дата и запросит данные для соответствующей стратегии
		err := chatbot.figiCallback(callbackQuery)
		if err != nil {
			chatbot.logger.Warn(err)
			return err
		}
		return nil
	case "new_worker":
		if chatbot.tradeBot == nil {
			resp := TGApi.NewMessage(chatId, "Сначала нужно ввести токен! Воспользуйтесь командой /tinkoff_token")
			_, err := chatbot.tg.Send(resp)
			return err
		}
		entry := chatbot.currChatState
		entry.isCallback = true
		entry.isCommand = false
		entry.value = callbackQuery.Data

		fmt.Print(callbackQuery.Data)
		//нужна функция которая возьмет выбор из дата и запросит данные для соответствующей стратегии
		err := chatbot.new_worker_callback(callbackQuery)
		if err != nil {
			chatbot.logger.Warn(err)
			return err
		}
		return nil
	case "delete_bot":
		entry := chatbot.currChatState
		entry.isCallback = true
		entry.isCommand = false
		entry.value = callbackQuery.Data
		chatbot.currChatState = entry
		err := chatbot.delete_botCallback(callbackQuery)
		if err != nil {
			chatbot.logger.Warn(err)
			return err
		}
		return nil
	default:
		return errors.New("unknown callback")
	}

}

var deleteKeyboard = TGApi.NewInlineKeyboardMarkup(TGApi.NewInlineKeyboardRow(
	TGApi.NewInlineKeyboardButtonData("удалить", "delete"),
	TGApi.NewInlineKeyboardButtonData("отмена", "cancel")))

func (chatbot *telegramChatBot) delete_bot_Command(msg *TGApi.Message) error {
	RespMsg := TGApi.NewMessage(msg.Chat.ID, "Вы действительно хотите удалить бота?")
	RespMsg.ReplyMarkup = deleteKeyboard
	_, err := chatbot.tg.Send(RespMsg)
	return err
}

func (chatbot *telegramChatBot) delete_botCallback(callback *TGApi.CallbackQuery) error {
	switch callback.Data {
	case "delete":
		if chatbot.tradeBot != nil {
			chatbot.StopTinkoffBot(callback.Message.Chat.ID)
			//delete(bot.signalsCh, callback.Message.Chat.ID)
			RespMsg := TGApi.NewMessage(callback.Message.Chat.ID, "Бот остановлен")
			_, err := chatbot.tg.Send(RespMsg)
			return err
		} else {
			RespMsg := TGApi.NewMessage(callback.Message.Chat.ID, "Бот остановлен")
			_, err := chatbot.tg.Send(RespMsg)
			return err
		}

	case "cancel":
		RespMsg := TGApi.NewMessage(callback.Message.Chat.ID, "Удаление бота отменено ")
		_, err := chatbot.tg.Send(RespMsg)
		return err
	default:
		return errors.New("unknown callback")
	}
}

func (chatbot *telegramChatBot) stop_worker_Command(msg *TGApi.Message) error {
	RespMsg := TGApi.NewMessage(msg.Chat.ID, "Введите id воркера, который необходимо остановить. "+
		"Чтобы получить список запущенных воркеров воспользуйтесь коммандой /worker_list")
	_, err := chatbot.tg.Send(RespMsg)
	return err
}

func (chatbot *telegramChatBot) stop_worker_Arg(msg *TGApi.Message) error {
	if chatbot.tradeBot.IsValidWorker(msg.Text) {
		id, err := ConvId(msg.Text)
		if err != nil {
			return err
		}
		chatbot.tradeBot.StopWorker(id)
		return nil
	} else {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Воркер с таким id: %s не найден. "+
			"Чтобы получить список запущенных воркеров воспользуйтесь коммандой /worker_list", msg.Text))
		_, err := chatbot.tg.Send(RespMsg)
		return err
	}
}

func (chatbot *telegramChatBot) worker_list_Command(msg *TGApi.Message) error {
	var infoSl []string
	var infoStr string
	infoList := chatbot.tradeBot.GetAllWorkersInfo()
	for _, v := range infoList {
		if len(infoStr)+len(v) >= 4095 {
			infoSl = append(infoSl, infoStr)
			infoStr = v
		} else {
			infoStr += fmt.Sprintf("%s\n\n", v)
		}
	}
	if len(infoSl) == 0 {
		infoSl = append(infoSl, infoStr)
	}
	RespMsg := TGApi.NewMessage(msg.Chat.ID, "Запущенные воркеры")
	_, err := chatbot.tg.Send(RespMsg)
	if err != nil {
		return err
	}
	for _, str := range infoSl {
		RespMsg = TGApi.NewMessage(msg.Chat.ID, str)
		_, err = chatbot.tg.Send(RespMsg)
		if err != nil {
			return err
		}
	}
	return nil
}

// START
func (chatbot *telegramChatBot) start_Command(msg *TGApi.Message) error {
	RespMsg := TGApi.NewMessage(msg.Chat.ID, "Привет! Это бот для Тинкофф Инвестиций. "+
		"Сначала необходимо получить в личном кабинете Тинькофф Инвестиций токен только для чтения. "+
		"Для передачи его боту воспользуйтесь командой /tinkoff_token")
	_, err := chatbot.tg.Send(RespMsg)
	return err
}

// TINKOFF_TOKEN
func (chatbot *telegramChatBot) tinkoff_token_Command(msg *TGApi.Message) error {
	RespMsg := TGApi.NewMessage(msg.Chat.ID, "Введите токен из Тинькофф Инвестиций")
	_, err := chatbot.tg.Send(RespMsg)
	return err
}

func (chatbot *telegramChatBot) tinkoff_token_Arg(msg *TGApi.Message) error {
	new_bot, err := trade.NewTradeBot(msg.Text, msg.Chat.ID)
	if err != nil {
		return err
	}
	if !new_bot.TokenIsValid() {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, "Неверный токен!")
		_, err = chatbot.tg.Send(RespMsg)
		return err
	}
	chatbot.tradeBot = new_bot
	//chatbot.botWg.Add(1)
	go chatbot.StartBot(msg.Chat.ID)
	RespMsg := TGApi.NewMessage(msg.Chat.ID, "Бот успешно создан!\n Чтобы отслеживать состояние актива "+
		"при помощи индикаторов воспользуйтесь командой /new_worker")
	_, err = chatbot.tg.Send(RespMsg)
	return err
}

// NEW_WORKER
var indicatorKeyboard = TGApi.NewInlineKeyboardMarkup(TGApi.NewInlineKeyboardRow(
	TGApi.NewInlineKeyboardButtonData("doubleEMA", "doubleEMA"),
	TGApi.NewInlineKeyboardButtonData("bollingerBands", "bollingerBands")))

func (chatbot *telegramChatBot) new_worker_Command(msg *TGApi.Message) error {
	/*buttons := TGApi.NewKeyboardButtonRow()
	buttons = append(buttons, TGApi.NewKeyboardButton("doubleEMA"))
	buttons = append(buttons, TGApi.NewKeyboardButton("BollingerBands"))
	RespMsg := TGApi.NewOneTimeReplyKeyboard(buttons)*/
	RespMsg := TGApi.NewMessage(msg.Chat.ID, "Выберите индикатор")
	RespMsg.ReplyMarkup = indicatorKeyboard
	_, err := chatbot.tg.Send(RespMsg)
	return err
}

// получаем выбраную стратегию и запрашиваем параметры для нее
func (chatbot *telegramChatBot) new_worker_callback(callback *TGApi.CallbackQuery) error {

	//надо как-то различать колл бэких
	var strat_fields []field
	switch callback.Data {
	case "bollingerBands":
		strat_fields = GetStrategyFieldsFromStruct(strategy.BollingerBands{})
	case "doubleEMA":
		strat_fields = GetStrategyFieldsFromStruct(strategy.DoubleEMA{})
	}

	var strat_fields_by_string string
	for k, v := range strat_fields {
		if k < len(strat_fields)-1 {
			strat_fields_by_string = strat_fields_by_string + fmt.Sprintf("%s %s, ", v.field_name, v.field_type)
		} else {
			strat_fields_by_string = strat_fields_by_string + fmt.Sprintf("%s %s", v.field_name, v.field_type)
		}
	}
	RespMsg := TGApi.NewMessage(callback.Message.Chat.ID,
		fmt.Sprintf("Введите через пробел figi актива и параметры для индикатора %s\n %s\n"+
			"Доступные для анализа интервалы: 1 час, 24 часа ",
			chatbot.currChatState.value, strat_fields_by_string))
	_, err := chatbot.tg.Send(RespMsg)
	return err
}

// FIGI
var figiKeyboard = TGApi.NewInlineKeyboardMarkup(TGApi.NewInlineKeyboardRow(
	TGApi.NewInlineKeyboardButtonData("bonds", "bonds"),
	TGApi.NewInlineKeyboardButtonData("etfs", "etfs"),
	TGApi.NewInlineKeyboardButtonData("shares", "shares")),
	TGApi.NewInlineKeyboardRow(TGApi.NewInlineKeyboardButtonData("currencies", "currencies"),
		TGApi.NewInlineKeyboardButtonData("futures", "futures")))

func (chatbot *telegramChatBot) get_figi_Command(msg *TGApi.Message) error {
	RespMsg := TGApi.NewMessage(msg.Chat.ID, "Выберите тип актива, для которого необходимо найти figi идентификатор")
	RespMsg.ReplyMarkup = figiKeyboard
	_, err := chatbot.tg.Send(RespMsg)
	return err
}

func (chatbot *telegramChatBot) get_figi_Arg(msg *TGApi.Message) error {
	keyWord := strings.Fields(msg.Text)

	assetsList, err := chatbot.tradeBot.GetAssetsList(chatbot.currChatState.value)
	var assetsStr string
	var tmpStr string
	var fullAssetsSl []string
	var keyAssetsSl []string
	for _, v := range assetsList {
		tmpStr = fmt.Sprintf("%s: %s, \n", v.Name, v.Figi)

		for _, s := range keyWord {
			if strings.Contains(strings.ToLower(v.Name), strings.ToLower(s)) {
				keyAssetsSl = append(keyAssetsSl, tmpStr)
				continue
			}
		}

		if len(assetsStr)+len(tmpStr) >= 4095 {
			fullAssetsSl = append(fullAssetsSl, assetsStr)
			assetsStr = tmpStr
		} else {
			assetsStr += tmpStr
		}
	}

	if len(keyAssetsSl) != 0 && len(keyAssetsSl) < 10 {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Найденные активы с таким названием"))
		_, err = chatbot.tg.Send(RespMsg)
		for _, p := range keyAssetsSl {
			RespMsg = TGApi.NewMessage(msg.Chat.ID, p)
			_, err = chatbot.tg.Send(RespMsg)
		}
	} else {
		RespMsg := TGApi.NewMessage(msg.Chat.ID, fmt.Sprintf("Активы с таким названием не найдены. Список всех доступных активов этого типа"))
		_, err = chatbot.tg.Send(RespMsg)
		for _, p := range fullAssetsSl {
			RespMsg = TGApi.NewMessage(msg.Chat.ID, p)
			_, err = chatbot.tg.Send(RespMsg)
		}
	}
	return err
}

// получаем выбраный тип актива и присылаем доступные для торговли
func (chatbot *telegramChatBot) figiCallback(callback *TGApi.CallbackQuery) error {
	RespMsg := TGApi.NewMessage(callback.Message.Chat.ID, "Введите название актива")
	_, err := chatbot.tg.Send(RespMsg)
	return err
}
