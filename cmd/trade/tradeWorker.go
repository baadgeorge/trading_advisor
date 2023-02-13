package trade

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"someshit/internal/strategy"
	"someshit/internal/strategy/utils"
	"someshit/pkg/proto"
	"someshit/pkg/sdk"
	"sync"
	"time"
)

//var services = sdk.NewServicePool()

type TradeWorker struct {
	workerID       uint32
	figi           string
	logger         *logrus.Entry
	workerSleepSec int
	strategy       strategy.Strategy
	indicatorState bool
	services       *sdk.ServicePool
	cancelCh       chan struct{}
}

// TODO sleep
func NewTradeWorker(workerConfig *WorkerConfig, services *sdk.ServicePool, logger *logrus.Entry) *TradeWorker {

	return &TradeWorker{
		workerID: workerConfig.workerId,
		figi:     workerConfig.figi,
		logger: logger.WithFields(logrus.Fields{
			"workerID": workerConfig.workerId,
			"figi":     workerConfig.figi,
		}),
		workerSleepSec: 15,
		strategy:       workerConfig.strategy,
		services:       services,
		indicatorState: false,
		cancelCh:       make(chan struct{}),
	}
}

func (tw *TradeWorker) Run(wg *sync.WaitGroup, stateChangesCh chan WorkersChanges) {
	defer wg.Done()
	stateChangesCh <- WorkersChanges{
		Img:         nil,
		WorkerId:    tw.workerID,
		SignalsType: Info_type,
		Description: fmt.Sprintf("Воркер с параметрами %s запущен", tw.GetWorkersDescr()),
	}
	var buf bytes.Buffer

	tw.logger.Infof("worker %d is running\n", tw.workerID)
	for {
		select {
		case <-time.After(time.Duration(tw.workerSleepSec) * time.Second):
			if !tw.tradingStatus() {
				continue
			}
			img, ind, err := tw.tradingIndicator()
			if err != nil {
				stateChangesCh <- WorkersChanges{
					Img:         nil,
					WorkerId:    tw.workerID,
					SignalsType: Err_type,
					Description: fmt.Sprintf("Ошибка в воркере %d: %s", tw.workerID, err),
				}
				tw.logger.Error(err)
				continue
			}
			// посылаем сигнал, только если индикатор изменился(не пустая строка)
			if ind != "" {
				if img != nil {
					_, err = img.WriteTo(&buf)
					if err != nil {
						stateChangesCh <- WorkersChanges{
							Img:         nil,
							WorkerId:    tw.workerID,
							SignalsType: Err_type,
							Description: fmt.Sprintf("Ошибка в воркере %d: %s", tw.workerID, err),
						}
						tw.logger.Error(err)
					}
				}
				stateChangesCh <- WorkersChanges{
					Img:         buf.Bytes(),
					WorkerId:    tw.workerID,
					SignalsType: Signal_type,
					Description: fmt.Sprintf("Получен новый сигнал от воркера с параметрами %s: %s",
						tw.GetWorkersDescr(), ind),
				}

			}
		case <-tw.cancelCh:
			stateChangesCh <- WorkersChanges{
				Img:         nil,
				WorkerId:    tw.workerID,
				SignalsType: Cancel_type,
				Description: fmt.Sprintf("Воркер %d остановлен", tw.workerID),
			}
			tw.logger.Infof("worker %d stopped\n", tw.workerID)
			return
		}
	}
}

func (tw *TradeWorker) tradingStatus() bool {
	status, err := tw.services.MarketDataService.GetTradingStatus(tw.figi)
	if err != nil {
		tw.logger.Errorf("can't get trading status: %v ", err)
		return false
	}
	tw.logger.Infof("trading status: %v ", status.TradingStatus.String())
	return status.TradingStatus == proto.SecurityTradingStatus_SECURITY_TRADING_STATUS_NORMAL_TRADING
}

const candlePartSize = 30

// доступность сигнала проверяется в Run()
func (tw *TradeWorker) tradingIndicator() (io.WriterTo, string, error) {
	interval := tw.strategy.GetAnalyzeInterval()
	if interval > 400 {
		tw.logger.Errorf("Too wide interval %d", interval)
		msg := fmt.Sprintf("Слишком большой интервал для анализа %d", interval)
		return nil, "", errors.New(msg)
	}

	var candles []*proto.HistoricCandle
	var rightTimePoint, leftTimePoint time.Time
	//var candlesPartCopy []*proto.HistoricCandle
	rightTimePoint = time.Now()

	for {
		leftTimePoint = rightTimePoint.Add(-time.Duration(candlePartSize) * time.Duration(tw.strategy.GetCandleInterval()) * time.Hour)
		candlesPart, err := tw.services.MarketDataService.GetCandles(tw.figi,
			leftTimePoint, rightTimePoint,
			time.Duration(tw.strategy.GetCandleInterval())*time.Hour)
		if err != nil {
			tw.logger.Error(err)
			msg := fmt.Sprint("Ошибка при получении свечей из API")
			return nil, "", errors.New(msg)
		}
		//copy(candlesPartCopy, candlesPart)
		if len(candles) == 0 {
			candles = append(candlesPart, candles...)
		} else if len(candlesPart) > 2 {
			candles = append(candlesPart[:len(candlesPart)-1], candles...)
		}

		/*	fmt.Println("\npart")
			for _, v := range candlesPart {
				fmt.Print(v.Time)
			}

			fmt.Println("\nfull")
			for _, v := range candles {
				fmt.Print(v.Time)
			}*/

		//нужно загрузить свечей в 2 раза больше, чем окно, для красивого графика =)
		if len(candles) >= 2*interval {
			break
		}
		rightTimePoint = leftTimePoint
	}

	ind, err := tw.strategy.Indicator(candles)
	if err != nil {
		tw.logger.Error(err)
		msg := "can't get indicator"
		return nil, "", errors.New(msg)
	}

	resp := ""
	var img io.WriterTo
	//если индикатор изменился
	if tw.indicatorState != ind {
		tw.indicatorState = ind //изменяем индикатор в экземпляре воркера
		if ind == true {
			resp = "ПОКУПКА"
		} else {
			resp = "ПРОДАЖА"
		}
		img, err = tw.strategy.DataPlot(utils.CandlesToTimeSeries(candles))
	}

	tw.logger.Infof("indiactor status, %t", ind)
	return img, resp, err
}

func (tw *TradeWorker) GetWorkersDescr() string {
	return fmt.Sprintf("worker id: %d strategy: %s", tw.workerID, tw.strategy.GetStrategyParamByString())
}

func (tw *TradeWorker) GetWorkerID() uint32 {
	return tw.workerID
}

func (tw *TradeWorker) GetWorkerCancelCh() chan struct{} {
	return tw.cancelCh
}
