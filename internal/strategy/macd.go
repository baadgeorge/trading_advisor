package strategy

import (
	"errors"
	"final/internal/strategy/utils"
	"final/pkg/proto"
	"fmt"
	"github.com/sdcoffey/techan"
)

type MACD struct {
	ShortWindow int
	LongWindow  int
	//SignalWindow        in
	CandleIntervalHours int
}

func (macd *MACD) Indicator(candles []*proto.HistoricCandle) (res IndicatorSignal, err error) {
	//нужен для обработки ошибок при вычислении индикатора
	defer func() {
		if r := recover(); r != nil {
			res = IndicatorSignal{
				Changed: false,
				Value:   false,
			}
			err = errors.New(fmt.Sprintf("panic in indicator func: %v", r))
		}
	}()

	convCandles := utils.CandlesToTimeSeries(candles)
	macdInd := techan.NewMACDIndicator(techan.NewClosePriceIndicator(convCandles), macd.ShortWindow, macd.LongWindow)
	macdCalc := macdInd.Calculate(len(convCandles.Candles) - 1).Float()
	closePrice := convCandles.LastCandle().ClosePrice.Float()

	//macd пробивает price сверху вниз - покупка
	if macdCalc > closePrice {
		return IndicatorSignal{
			Changed: true,
			Value:   true,
		}, nil
	}

	//macd пробивает price сверху вниз - покупка
	if macdCalc < closePrice {
		return IndicatorSignal{
			Changed: true,
			Value:   true,
		}, nil
	}

	//macd пробивает price снизу вверх - продажа
	if macdCalc > closePrice {
		return IndicatorSignal{
			Changed: true,
			Value:   false,
		}, nil
	}

	return IndicatorSignal{
		Changed: false,
		Value:   false,
	}, nil
}

func (macd *MACD) GetCandleInterval() int {
	return macd.CandleIntervalHours
}

func (macd *MACD) GetAnalyzeInterval() int {
	if macd.LongWindow > macd.ShortWindow {
		return macd.LongWindow
	}
	return macd.ShortWindow
}

func (macd *MACD) GetStrategyParamByString() string {
	return fmt.Sprintf("MACD: ShortWindow: %d LongWindow: %d CandleIntervalHours: %d", macd.ShortWindow, macd.LongWindow, macd.CandleIntervalHours)
}

func (macd *MACD) DataPlot(convCandles *techan.TimeSeries) ([][]byte, error) {
	var candleSl []utils.PlotItemStruct
	var macdSl []utils.PlotItemStruct
	var dataPlots [][]byte

	for k, v := range convCandles.Candles {
		if k < macd.GetAnalyzeInterval() {
			continue
		}

		//candlePart = utils.CandlesToTimeSeries(candles[:k])
		candleSl = append(candleSl, utils.PlotItemStruct{
			Period: v.Period,
			Value:  v.ClosePrice.Float(),
		})
		macdSl = append(macdSl, utils.PlotItemStruct{
			Period: v.Period,
			Value:  techan.NewMACDIndicator(techan.NewClosePriceIndicator(convCandles), macd.ShortWindow, macd.LongWindow).Calculate(k).Float(),
		})

	}

	fmt.Printf("candles\n %v\n", candleSl)
	fmt.Printf("macd\n %v\n", macdSl)

	candleMap := make(map[string][]utils.PlotItemStruct)
	candleMap["candles"] = candleSl
	cp, err := utils.PlotData(candleMap, "time", "price", "Close price")
	if err != nil {
		return nil, err
	}
	dataPlots = append(dataPlots, cp)

	indPlot := make(map[string][]utils.PlotItemStruct)
	indPlot["macd"] = macdSl
	ip, err := utils.PlotData(indPlot, "time", "value", "MACD")
	if err != nil {
		return nil, err
	}
	dataPlots = append(dataPlots, ip)

	return dataPlots, nil
}
