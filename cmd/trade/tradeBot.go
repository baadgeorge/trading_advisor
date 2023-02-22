package trade

import (
	"errors"
	"github.com/sirupsen/logrus"
	"someshit/pkg/proto"
	"someshit/pkg/sdk"
	"strconv"
	"sync"
	"time"
)

type WorkersChanges struct {
	Img         []byte
	WorkerId    uint32
	SignalsType State
	Description string
}

// bot for 1 tinkoff account with multiple goroutines
// new workers configuration get from newWorkerConfigCh
// running workers added to waitGroup workersWg
// workers id and its description contain in workersInfoList
// workers cancel channel contain in workersCancelChannels
type TradeBot struct {
	token             string
	botCloseCh        chan struct{}
	workersChangesCh  chan WorkersChanges
	newWorkerConfigCh chan *WorkerConfig
	workersWg         *sync.WaitGroup
	workers           map[uint32]*TradeWorker
	sdkServices       *sdk.ServicePool
	logger            *logrus.Entry
	accountID         int64
}

func NewTradeBot(token string, accountID int64) (*TradeBot, error) {

	serv, err := sdk.NewServicePool(token)
	if err != nil {
		return nil, err
	}
	return &TradeBot{
		token:             token,
		botCloseCh:        make(chan struct{}),
		workersChangesCh:  make(chan WorkersChanges),
		newWorkerConfigCh: make(chan *WorkerConfig),
		workersWg:         new(sync.WaitGroup),
		//ctx:             ctx,
		workers:     make(map[uint32]*TradeWorker),
		sdkServices: serv,
		logger: logrus.WithFields(logrus.Fields{
			"time":      time.Now(),
			"accountID": accountID,
		}),
		accountID: accountID,
	}, nil
}

func (tb *TradeBot) StopBot() {
	tb.botCloseCh <- struct{}{}
	close(tb.botCloseCh)
	return
}

// нужно создавать воркеров здесь, чтобы передавать им один экз сервисов и токен
func (tb *TradeBot) StartNewWorkers() {
	///\
	//defer wg.Done()
	for {
		select {
		case config := <-tb.newWorkerConfigCh:
			worker := NewTradeWorker(config, tb.sdkServices, tb.logger)
			tb.workers[worker.workerID] = worker
			tb.workersWg.Add(1)
			go worker.Run(tb.workersWg, tb.workersChangesCh)
		case <-tb.botCloseCh:
			for k := range tb.workers {
				tb.StopWorker(k)
			}
			return
		}
	}
}
func (tb *TradeBot) StopWorker(workerID uint32) {
	//doesn't need wg done()
	//worker kills by defer done() in Run func
	//defer tb.workersWg.Done()
	tb.workers[workerID].GetWorkerCancelCh() <- struct{}{}
	close(tb.workers[workerID].GetWorkerCancelCh())
	delete(tb.workers, workerID)
	return
}

func (tb *TradeBot) GetAllWorkersInfo() []string {
	var info []string
	for _, v := range tb.workers {
		info = append(info, v.GetWorkersDescr())
	}
	return info
}

func (tb *TradeBot) IsValidWorker(id string) bool {
	convId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return false
	}
	_, ok := tb.workers[uint32(convId)]
	return ok
}

func (tb *TradeBot) GetNewWorkersConfigCh() chan *WorkerConfig {
	return tb.newWorkerConfigCh
}

func (tb *TradeBot) GetWorkersChangesCh() chan WorkersChanges {
	return tb.workersChangesCh
}

type assets struct {
	Name string
	Figi string
}

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
			//fmt.Printf("Name: %s, figi: %s, currency: %s\n", v.Name, v.figi, v.Nominal.Currency)
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

func (tb *TradeBot) TokenIsValid() bool {
	_, err := tb.sdkServices.UsersService.GetInfo()
	if err != nil {
		return false
	}
	return true
}
