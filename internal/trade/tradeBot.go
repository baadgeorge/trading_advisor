package trade

import (
	"errors"
	"final/pkg/proto"
	"final/pkg/sdk"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// структура для передачи сигналов из воркера
type WorkersChanges struct {
	Img         [][]byte
	WorkerId    uint32
	SignalsType State
	Description string
}

// структура для сигнального бота, который привязан к 1 аккаунту
// конфигурации новых воркеров приходят через канал newWorkerConfigCh
// запущенные воркеры находятся в workers в соответствии с их идентификаторами
type TradeBot struct {
	//token             string                  //токен, который необходим для доступа к сервисам тинькофф инвестиций
	botCloseCh        chan struct{}           //канал для остановки бота и всех его воркеров
	workersChangesCh  chan WorkersChanges     //канал для получения сигналов(в т.ч. ошибок) от своих воркеров
	newWorkerConfigCh chan *WorkerConfig      //канал для получения конфигурации воркеров, которые необходимо запустить
	workers           map[uint32]*TradeWorker //запущенные экземпляры воркеров
	sdkServices       *sdk.ServicePool        //сервисы тинькофф инвестиций
	logger            *logrus.Logger
	accountID         int64 //уникальный идентификатор бота
}

// функция создания нового экземпляра сигнального бота
func NewTradeBot(accountID int64, token string) (*TradeBot, error) {

	serv, err := sdk.NewServicePool(token)
	if err != nil {
		return nil, err
	}

	log := logrus.New()
	log.SetReportCaller(true)
	log.WithFields(logrus.Fields{
		"time":      time.Now(),
		"accountID": accountID,
	})
	return &TradeBot{
		//token:             token,
		botCloseCh:        make(chan struct{}),
		workersChangesCh:  make(chan WorkersChanges),
		newWorkerConfigCh: make(chan *WorkerConfig),
		//workersWg:         new(sync.WaitGroup),
		//ctx:             ctx,
		workers:     make(map[uint32]*TradeWorker),
		sdkServices: serv,
		logger:      log,
		accountID:   accountID,
	}, nil
}

// метод остановки бота и всех его воркеров
func (tb *TradeBot) StopBot() {
	tb.botCloseCh <- struct{}{}
	close(tb.botCloseCh)
	return
}

// метод запуска и остановки воркеров бота
func (tb *TradeBot) TradeBotRun() {
	for {
		select {
		//создание новых воркеров в боте
		case config := <-tb.newWorkerConfigCh:
			// нужно создавать воркеров здесь, чтобы передавать им один экземпляр сервисов и токен
			worker := NewTradeWorker(config, tb.sdkServices, tb.logger)
			tb.workers[worker.workerID] = worker
			go worker.Run(tb.workersChangesCh)

			//остановка всех воркеров бота
		case <-tb.botCloseCh:
			for k := range tb.workers {
				tb.StopWorker(k)
			}
			return
		}
	}
}

// метод завершения воркера с заданным id
func (tb *TradeBot) StopWorker(workerID uint32) {
	//worker kills by defer done() in Run func
	tb.workers[workerID].GetWorkerCancelCh() <- struct{}{}
	close(tb.workers[workerID].GetWorkerCancelCh())
	delete(tb.workers, workerID)
	return
}

// метод получения информации о всех запущенных воркерах
func (tb *TradeBot) GetAllWorkersInfo() []string {
	var info []string
	for _, v := range tb.workers {
		info = append(info, v.GetWorkersDescr())
	}
	return info
}

// метод проверки существования вокркера с заданным id
func (tb *TradeBot) IsValidWorker(id string) bool {
	convId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return false
	}
	_, ok := tb.workers[uint32(convId)]
	return ok
}

// метод возвращающий канал получения конфигураций новых воркеров
func (tb *TradeBot) GetNewWorkersConfigCh() chan *WorkerConfig {
	return tb.newWorkerConfigCh
}

// метод возвращающий канал получения сообщений от воркеров
func (tb *TradeBot) GetWorkersChangesCh() chan WorkersChanges {
	return tb.workersChangesCh
}

// структура актива
type assets struct {
	Name string
	Figi string
}

// метод получения списка активовов заданного типа
func (tb *TradeBot) GetAssetsList(instrumentsType string) ([]assets, error) {
	var assetsList []assets
	switch instrumentsType {
	case "bonds":
		bonds, err := tb.sdkServices.InstrumentsService.Bonds(proto.InstrumentStatus_INSTRUMENT_STATUS_BASE)
		if err != nil {
			return nil, err
		}
		for _, v := range bonds {
			assetsList = append(assetsList, assets{
				Name: v.Name,
				Figi: v.Figi,
			})
			//fmt.Printf("Name: %s, figi: %s, currency: %s\n", v.Name, v.figi, v.Nominal.Currency)
		}
		return assetsList, nil
	case "currencies":
		currencies, err := tb.sdkServices.InstrumentsService.Currencies(proto.InstrumentStatus_INSTRUMENT_STATUS_BASE)
		if err != nil {
			return nil, err
		}
		for _, v := range currencies {
			assetsList = append(assetsList, assets{
				Name: v.Name,
				Figi: v.Figi,
			})
			fmt.Printf("Name: %s, figi: %s, currency: %s\n", v.Name, v.Figi, v.Nominal.Currency)
		}
		return assetsList, nil
	case "etfs":
		etfs, err := tb.sdkServices.InstrumentsService.Etfs(proto.InstrumentStatus_INSTRUMENT_STATUS_BASE)
		if err != nil {
			return nil, err
		}
		for _, v := range etfs {
			assetsList = append(assetsList, assets{
				Name: v.Name,
				Figi: v.Figi,
			})
			//fmt.Printf("Name: %s, figi: %s, currency: %s\n", v.Name, v.figi, v.Currency)
		}
		return assetsList, nil
	case "futures":
		futures, err := tb.sdkServices.InstrumentsService.Futures(proto.InstrumentStatus_INSTRUMENT_STATUS_BASE)
		if err != nil {
			return nil, err
		}
		for _, v := range futures {
			assetsList = append(assetsList, assets{
				Name: v.Name,
				Figi: v.Figi,
			})
			//fmt.Printf("Name: %s, figi: %s, currency: %s\n", v.Name, v.figi, v.Currency)
		}
		return assetsList, nil
	case "shares":
		shares, err := tb.sdkServices.InstrumentsService.Shares(proto.InstrumentStatus_INSTRUMENT_STATUS_BASE)
		if err != nil {
			return nil, err
		}
		for _, v := range shares {
			assetsList = append(assetsList, assets{
				Name: v.Name,
				Figi: v.Figi,
			})
			//fmt.Printf("Name: %s, figi: %s, currency: %s\n", v.Name, v.figi, v.Currency)
		}
		return assetsList, nil
	default:
		return nil, errors.New("unknown type of assets")
	}
}

// метод проверки токена на валидность
func (tb *TradeBot) TokenIsValid() bool {
	_, err := tb.sdkServices.UsersService.GetInfo()
	if err != nil {
		return false
	}
	return true
}
