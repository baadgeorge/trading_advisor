package strategy

import (
	"errors"
	"fmt"

	"final/internal/strategy/utils"
	"final/pkg/proto"
	"github.com/sdcoffey/techan"
)

type RSI struct {
	HighBorder          float64
	LowBorder           float64
	CandleIntervalHours int
	Window              int
	Position            int `reflect:"-"`
}

const (
	RSIaboveHighBorder = 1
	RSIbetweenBorders  = 2
	RSIunderLowBorder  = 3
)

func (rsi *RSI) Indicator(candles []*proto.HistoricCandle) (res IndicatorSignal, err error) {
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
	rsiInd := techan.NewRelativeStrengthIndexIndicator(techan.NewClosePriceIndicator(convCandles), rsi.Window).Calculate(rsi.Window).Float()

	//rsi пробивает high border сверху вниз - продажа
	if rsi.Position == RSIaboveHighBorder && rsiInd <= rsi.HighBorder {
		rsi.Position = RSIbetweenBorders
		return IndicatorSignal{
			Changed: true,
			Value:   true,
		}, nil
	}

	//rsi пробивает low border снизу вверх - покупка
	if rsi.Position == RSIunderLowBorder && rsiInd >= rsi.LowBorder {
		rsi.Position = RSIbetweenBorders
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

func (rsi *RSI) GetCandleInterval() int {
	return rsi.CandleIntervalHours
}

func (rsi *RSI) GetAnalyzeInterval() int {
	return rsi.Window
}

func (rsi *RSI) GetStrategyParamByString() string {
	return fmt.Sprintf("RSI: HighBorder: %f LowBorder: %f Window: %d CandleIntervalHours: %d",
		rsi.HighBorder, rsi.LowBorder, rsi.Window, rsi.CandleIntervalHours)
}

func (rsi *RSI) DataPlot(convCandles *techan.TimeSeries) ([][]byte, error) {
	var candleSl []utils.PlotItemStruct
	var rsiSl []utils.PlotItemStruct
	var dataPlots [][]byte

	//for k, v := range convCandles.Candles {
	candles := convCandles.Candles

	for k := rsi.GetAnalyzeInterval() - 1; k < len(candles); k++ {
		/*	if k < macd.GetAnalyzeInterval() {
			continue
		}*/

		//candlePart = utils.CandlesToTimeSeries(candles[:k])
		candleSl = append(candleSl, utils.PlotItemStruct{
			Period: candles[k].Period,
			Value:  candles[k].ClosePrice.Float(),
		})
		rsiSl = append(rsiSl, utils.PlotItemStruct{
			Period: candles[k].Period,
			Value:  techan.NewRelativeStrengthIndexIndicator(techan.NewClosePriceIndicator(convCandles), rsi.Window).Calculate(k).Float(),
		})

	}

	fmt.Printf("candles\n %v\n", candleSl)
	fmt.Printf("rsi\n %v\n", rsiSl)

	candleMap := make(map[string][]utils.PlotItemStruct)
	candleMap["candles"] = candleSl
	cp, err := utils.PlotData(candleMap, "time", "price", "Close price")
	if err != nil {
		return nil, err
	}
	dataPlots = append(dataPlots, cp)

	indPlot := make(map[string][]utils.PlotItemStruct)
	indPlot["rsi"] = rsiSl
	ip, err := utils.PlotData(indPlot, "time", "value", "RSI")
	if err != nil {
		return nil, err
	}
	dataPlots = append(dataPlots, ip)

	return dataPlots, nil
}
