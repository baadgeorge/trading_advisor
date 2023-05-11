package trade

import (
	"errors"
	"final/internal/strategy"
	"final/internal/strategy/utils"
	"final/pkg/proto"
	"final/pkg/sdk"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// структура воркера
type TradeWorker struct {
	workerID       uint32 //уникальный идентификатор воркера
	figi           string //figi актива
	assetsName     string //название актива
	logger         *logrus.Entry
	workerSleepSec int                     //интервал запроса данных
	strategy       strategy.Strategy       //стратегия воркера(индикатор)
	indicatorState bool                    //текущее состояние индикатора
	services       *sdk.ServicePool        //сервисы Тинькофф Инвестиций
	candles        []*proto.HistoricCandle //свечи
	cancelCh       chan struct{}           //канал завершения воркера
}

// функция создания экземпляра воркера
func NewTradeWorker(workerConfig *WorkerConfig, services *sdk.ServicePool, logger *logrus.Logger) *TradeWorker {
	return &TradeWorker{
		workerID:   workerConfig.workerId,
		figi:       workerConfig.figi,
		assetsName: "",
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

// метод инициализации данных воркера(свечи и название актива)
func (tw *TradeWorker) initData(stateChangesCh chan WorkersChanges) error {
	name, err := tw.GetAssetsNameByFigi()
	if err != nil {
		stateChangesCh <- WorkersChanges{
			Img:         nil,
			WorkerId:    tw.workerID,
			SignalsType: Err_type,
			Description: fmt.Sprintf("Ошибка в воркере %d: неудалось загрузить имя актива", tw.workerID),
		}
		tw.logger.Error(err)
		return err
	}
	tw.assetsName = name

	err = tw.initCandles()
	if err != nil {
		stateChangesCh <- WorkersChanges{
			Img:         nil,
			WorkerId:    tw.workerID,
			SignalsType: Err_type,
			Description: fmt.Sprintf("Ошибка в воркере %d: %v", tw.workerID, err),
		}
		tw.logger.Error(err)
		return err
	}
	return nil
}

// метод циклической проверки индикатора и остановки воркера
func (tw *TradeWorker) Run(stateChangesCh chan WorkersChanges) {
	//defer wg.Done()
	count := 0
	tw.logger.Infof("worker %d is running\n", tw.workerID)
	stateChangesCh <- WorkersChanges{
		Img:         nil,
		WorkerId:    tw.workerID,
		SignalsType: Info_type,
		Description: fmt.Sprintf("Воркер с параметрами %s запущен", tw.GetWorkersDescr()),
	}
	for {
		select {
		case <-time.After(time.Duration(tw.workerSleepSec) * time.Second):

			if count == 0 {
				err := tw.initData(stateChangesCh)
				count++
				//если есть ошибка, то переходим на след. интерацию цикла, чтобы ее обработать
				if err != nil {
					continue
				}
			}

			//TODO !!!
			/*if !tw.tradingStatus() {
				continue
			}*/
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
				//var bytesSl [][]byte

				//buf := make([]bytes.Buffer, len(img))
				if img != nil {
					//var buf bytes.Buffer
					/*for i, v := range img {
						_, err = v.WriteTo(&buf[i])
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

						bytesSl = append(bytesSl, buf[i].Bytes())
					}*/
					stateChangesCh <- WorkersChanges{
						Img:         img,
						WorkerId:    tw.workerID,
						SignalsType: Signal_type,
						Description: fmt.Sprintf("Получен новый сигнал от воркера с параметрами %s\n %v",
							tw.GetWorkersDescr(), ind),
					}
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

// метод добавления новых свечей в candles
func (tw *TradeWorker) addNewCandles() error {
	if len(tw.candles) == 0 {
		return errors.New("Historic candles were not init")
	}
	//fmt.Printf("last %v\n", tw.candles[len(tw.candles)-1].GetTime().AsTime().Local())
	//fmt.Printf("now %v\n", time.Now())
	lastCandle := tw.candles[len(tw.candles)-1].GetTime().AsTime().Local()
	if lastCandle.Sub(time.Now()) < time.Duration(tw.strategy.GetCandleInterval())*time.Hour {
		return nil
	}
	candles, err := tw.services.MarketDataService.GetCandles(tw.figi,
		lastCandle, time.Now(),
		time.Duration(tw.strategy.GetCandleInterval())*time.Hour)
	if err != nil {
		return err
	}
	if len(candles) == 0 {
		return nil
	}
	if len(candles) == 1 && !candles[0].IsComplete {
		return nil
	}
	//если последняя свеча не закрыта, то не добавляем её
	if !candles[len(candles)-1].IsComplete {
		tw.candles = append(tw.candles[len(candles)-1:], candles[:len(candles)-1]...)
		return nil
	}
	tw.candles = append(tw.candles[len(candles):], candles...)
	return nil
}

const candlePartSize = 23

// метод инициализации свечей
func (tw *TradeWorker) initCandles() error {
	interval := tw.strategy.GetAnalyzeInterval()
	if interval > 400 {
		tw.logger.Errorf("Too wide interval %d", interval)
		msg := fmt.Sprintf("Слишком большой интервал для анализа %d", interval)
		return errors.New(msg)
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
			return errors.New(msg)
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
		if len(candles) > 2*interval {
			break
		}
		rightTimePoint = leftTimePoint
	}
	tw.candles = candles
	return nil
}

// метод проверки доступности торговли активом
func (tw *TradeWorker) tradingStatus() bool {
	status, err := tw.services.MarketDataService.GetTradingStatus(tw.figi)
	if err != nil {
		tw.logger.Errorf("can't get trading status: %v ", err)
		return false
	}
	tw.logger.Infof("trading status: %v ", status.TradingStatus.String())
	return status.TradingStatus == proto.SecurityTradingStatus_SECURITY_TRADING_STATUS_NORMAL_TRADING
}

// метод проверки индикатора
func (tw *TradeWorker) tradingIndicator() ([][]byte, string, error) {
	// доступность сигнала проверяется в Run()
	err := tw.addNewCandles()
	if err != nil {
		tw.logger.Error(err)
		msg := "ошибка при добавлении новых свечей"
		return nil, "", errors.New(msg)
	}

	ind, err := tw.strategy.Indicator(tw.candles)
	if err != nil {
		tw.logger.Error(err)
		msg := "ошибка при проверке индикатора"
		return nil, "", errors.New(msg)
	}

	resp := ""
	var img [][]byte

	//TODO !!!!!
	//если индикатор изменился
	//if tw.indicatorState != ind.Value && ind.Changed == true {
	if tw.indicatorState == false {
		tw.indicatorState = ind.Value //изменяем индикатор в экземпляре воркера
		if ind.Value == true {
			resp = "ПОКУПКА"
		} else {
			resp = "ПРОДАЖА"
		}

		img, err = tw.strategy.DataPlot(utils.CandlesToTimeSeries(tw.candles))
	}

	tw.logger.Infof("indiactor status, %t", ind)
	return img, resp, err
}

// метод получения описания воркера
func (tw *TradeWorker) GetWorkersDescr() string {
	return fmt.Sprintf("worker id: %d figi: %s, assets name: %s, strategy %s", tw.workerID, tw.figi, tw.assetsName, tw.strategy.GetStrategyParamByString())
}

// метод получения id воркера
func (tw *TradeWorker) GetWorkerID() uint32 {
	return tw.workerID
}

// метод получения канала остановки вокера
func (tw *TradeWorker) GetWorkerCancelCh() chan struct{} {
	return tw.cancelCh
}

// метод получения названия актива по его figi
func (tw *TradeWorker) GetAssetsNameByFigi() (string, error) {

	tmp, err := tw.services.InstrumentsService.GetInstrumentBy(
		proto.InstrumentRequest{
			IdType: proto.InstrumentIdType_INSTRUMENT_ID_TYPE_FIGI,
			Id:     tw.figi,
		})
	if err != nil {
		return "", err
	}
	return tmp.Name, nil
}
