package strategy

import (
	"errors"
	"fmt"
	"github.com/sdcoffey/techan"
	"someshit/internal/strategy/utils"
	"someshit/pkg/proto"
)

type BollingerBands struct {
	//AnalyzeIntervalInHours int
	Window              int
	Sigma               float64 //отклонение(влияет на ширину канала)
	CandleIntervalHours int
	//PriceBetweenBands bool
}

func NewBollingerBand(candleIntervalHours, window int, sigma float64) *BollingerBands {
	return &BollingerBands{
		CandleIntervalHours: candleIntervalHours,
		//AnalyzeIntervalInHours: analyzeInterval,
		Window: window,
		Sigma:  sigma,
		//PriceBetweenBands: true,
	}
}

// при пробитии ценой нижней границы - покупка, при пробитии верхней границы - продажа
func (bb *BollingerBands) Indicator(candles []*proto.HistoricCandle) (res IndicatorSignal, err error) {
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

	lowerBBI := techan.NewBollingerLowerBandIndicator(techan.NewClosePriceIndicator(convCandles), bb.Window, bb.Sigma)
	upperBBI := techan.NewBollingerUpperBandIndicator(techan.NewClosePriceIndicator(convCandles), bb.Window, bb.Sigma)
	calcLower := lowerBBI.Calculate(len(convCandles.Candles) - 1).Float()
	calcUpper := upperBBI.Calculate(len(convCandles.Candles) - 1).Float()
	closePrice := convCandles.LastCandle().ClosePrice.Float()

	//при пробитии ценой нижней границы - покупка
	if closePrice <= calcLower {
		return IndicatorSignal{
			Changed: true,
			Value:   true,
		}, nil
	}

	//при пробитии ценой верхней границы - продажа
	if closePrice >= calcUpper {
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

func (bb *BollingerBands) GetCandleInterval() int {
	return bb.CandleIntervalHours
}

func (bb *BollingerBands) GetAnalyzeInterval() int {
	return bb.Window

}

func (bb *BollingerBands) GetStrategyParamByString() string {
	return fmt.Sprintf("Bollinger Bands: Window: %d Sigma %f CandleIntervalHours: %d\n", bb.Window, bb.Sigma, bb.CandleIntervalHours)
}

func (bb *BollingerBands) DataPlot(convCandles *techan.TimeSeries) ([][]byte, error) {
	var candleSl []utils.PlotItemStruct
	var upperSl []utils.PlotItemStruct
	var lowerSl []utils.PlotItemStruct
	var dataPlots [][]byte
	/*candleSl := make(map[techan.TimePeriod]float64)
	shortSl := make(map[techan.TimePeriod]float64)
	longSl := make(map[techan.TimePeriod]float64)*/

	for k, v := range convCandles.Candles {
		if k < bb.GetAnalyzeInterval() {
			continue
		}

		//candlePart = utils.CandlesToTimeSeries(candles[:k])
		candleSl = append(candleSl, utils.PlotItemStruct{
			Period: v.Period,
			Value:  v.ClosePrice.Float(),
		})
		lowerSl = append(lowerSl, utils.PlotItemStruct{
			Period: v.Period,
			Value:  techan.NewBollingerLowerBandIndicator(techan.NewClosePriceIndicator(convCandles), bb.Window, bb.Sigma).Calculate(k).Float(),
		})
		upperSl = append(upperSl, utils.PlotItemStruct{
			Period: v.Period,
			Value:  techan.NewBollingerUpperBandIndicator(techan.NewClosePriceIndicator(convCandles), bb.Window, bb.Sigma).Calculate(k).Float(),
		})
		/*[v.Period] = v.ClosePrice.Float()
		shortSl[v.Period] = techan.NewEMAIndicator(techan.NewClosePriceIndicator(convCandles), dema.ShortWindow).Calculate(k).Float()
		longSl[v.Period] = techan.NewEMAIndicator(techan.NewClosePriceIndicator(convCandles), dema.LongWindow).Calculate(k).Float()*/
	}

	fmt.Printf("candles\n %v\n", candleSl)
	fmt.Printf("lower\n %v\n", lowerSl)
	fmt.Printf("upper\n %v\n", upperSl)

	candleMap := make(map[string][]utils.PlotItemStruct)
	candleMap["candles"] = candleSl
	candleMap["upperEMA"] = upperSl
	candleMap["lowerEMA"] = lowerSl
	p, err := utils.PlotData(candleMap, "time", "price", "Bollinger Bands")
	if err != nil {
		return nil, err
	}

	/*k, err := p.WriterTo(30*vg.Centimeter, 15*vg.Centimeter, "png")
	if err != nil {
		return nil, err
	}*/
	dataPlots = append(dataPlots, p)

	return dataPlots, nil
}
