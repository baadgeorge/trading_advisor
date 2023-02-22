package strategy

import (
	"errors"
	"fmt"
	"github.com/sdcoffey/techan"
	"gonum.org/v1/plot/vg"
	"io"
	"someshit/internal/strategy/utils"
	"someshit/pkg/proto"
)

const (
	longEMA  = true
	shortEMA = false
)

type DoubleEMA struct {
	ShortWindow         int
	LongWindow          int
	CandleIntervalHours int
	//AnalyzeIntervalHours int
	//OffsetShort          int
	//OffsetLong             int
	WhichEMAHigher bool `reflect:"-"`
}

func NewDoubleEMA(shortWindow, longWindow, candleInterval, offsetShort, offsetLong int) *DoubleEMA {
	return &DoubleEMA{
		ShortWindow:         shortWindow,
		LongWindow:          longWindow,
		CandleIntervalHours: candleInterval,
		//AnalyzeIntervalHours: analyzeInterval,
		//OffsetShort:            offsetShort,
		//OffsetLong:             offsetLong,
		WhichEMAHigher: false,
	}
}

// если longEMA пробивает shortEMA снизу вверх - продажа(false), сверху вниз - покупка(true)
func (dema *DoubleEMA) Indicator(candles []*proto.HistoricCandle) (res IndicatorSignal, err error) {
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
	short := techan.NewEMAIndicator(techan.NewClosePriceIndicator(convCandles), dema.ShortWindow)
	long := techan.NewEMAIndicator(techan.NewClosePriceIndicator(convCandles), dema.LongWindow)

	shortCalc := short.Calculate(len(convCandles.Candles) - 1).Float()
	longCalc := long.Calculate(len(convCandles.Candles) - 1).Float()

	//shortEMA пробивает longEMA сверху вниз
	if shortCalc > longCalc && dema.WhichEMAHigher == longEMA {
		dema.WhichEMAHigher = shortEMA
		return IndicatorSignal{
			Changed: true,
			Value:   true,
		}, nil
	}
	/*//longEMA ниже shortEMA и положение не изменилось
	if shortCalc >= longCalc && dema.WhichEMAHigher == shortEMA {
		return true, nil
	}*/
	//longEMA пробивает shortEMA снизу вверх
	if longCalc > shortCalc && dema.WhichEMAHigher == shortEMA {
		dema.WhichEMAHigher = longEMA
		return IndicatorSignal{
			Changed: true,
			Value:   false,
		}, nil
	}
	/*//shortEMA ниже longEMA и положение не изменилось
	if longCalc >= shortCalc && dema.WhichEMAHigher == longEMA {
		return false, nil
	}
	*/
	return IndicatorSignal{
		Changed: false,
		Value:   false,
	}, nil
}

func (dema *DoubleEMA) DataPlot(convCandles *techan.TimeSeries) (io.WriterTo, error) {
	var candleSl []utils.CandleStruct
	var shortSl []utils.CandleStruct
	var longSl []utils.CandleStruct
	/*candleSl := make(map[techan.TimePeriod]float64)
	shortSl := make(map[techan.TimePeriod]float64)
	longSl := make(map[techan.TimePeriod]float64)*/

	for k, v := range convCandles.Candles {
		if k < dema.GetAnalyzeInterval() {
			continue
		}

		//candlePart = utils.CandlesToTimeSeries(candles[:k])
		candleSl = append(candleSl, utils.CandleStruct{
			Period: v.Period,
			Value:  v.ClosePrice.Float(),
		})
		shortSl = append(shortSl, utils.CandleStruct{
			Period: v.Period,
			Value:  techan.NewEMAIndicator(techan.NewClosePriceIndicator(convCandles), dema.ShortWindow).Calculate(k).Float(),
		})
		longSl = append(longSl, utils.CandleStruct{
			Period: v.Period,
			Value:  techan.NewEMAIndicator(techan.NewClosePriceIndicator(convCandles), dema.LongWindow).Calculate(k).Float(),
		})
		/*[v.Period] = v.ClosePrice.Float()
		shortSl[v.Period] = techan.NewEMAIndicator(techan.NewClosePriceIndicator(convCandles), dema.ShortWindow).Calculate(k).Float()
		longSl[v.Period] = techan.NewEMAIndicator(techan.NewClosePriceIndicator(convCandles), dema.LongWindow).Calculate(k).Float()*/
	}

	fmt.Printf("candles\n %v\n", candleSl)
	fmt.Printf("short\n %v\n", shortSl)
	fmt.Printf("long\n %v\n", longSl)

	candleMap := make(map[string][]utils.CandleStruct)
	candleMap["candles"] = candleSl
	candleMap["shortEMA"] = shortSl
	candleMap["longEMA"] = longSl
	p, err := utils.CandlesToPlot(candleMap)
	if err != nil {
		return nil, err
	}

	return p.WriterTo(30*vg.Centimeter, 15*vg.Centimeter, "png")

}

func (dema *DoubleEMA) GetCandleInterval() int {
	return dema.CandleIntervalHours
}
func (dema *DoubleEMA) GetAnalyzeInterval() int {
	if dema.LongWindow >= dema.ShortWindow {
		return dema.LongWindow
	}
	return dema.ShortWindow
}

func (dema *DoubleEMA) GetStrategyParamByString() string {
	return fmt.Sprintf("Double EMA: ShortWindow: %d LongWindow: %d CandleIntervalHours: %d", dema.ShortWindow, dema.LongWindow, dema.CandleIntervalHours)
}
