package strategy

import (
	"errors"
	"fmt"
	"github.com/sdcoffey/techan"
	"someshit/internal/strategy/utils"
	"someshit/pkg/proto"
)

type BollingerBands struct {
	CandleIntervalInHours  int
	AnalyzeIntervalInHours int
	Window                 int
	Sigma                  float64 //отклонение(влияет на ширину канала)
}

func NewBollingerBand(candleIntervalHours, analyzeInterval, window int, sigma float64) *BollingerBands {
	return &BollingerBands{
		CandleIntervalInHours:  candleIntervalHours,
		AnalyzeIntervalInHours: analyzeInterval,
		Window:                 window,
		Sigma:                  sigma,
	}
}

// при пробитии ценой нижней границы - покупка, при пробитии верхней границы - продажа
func (bb *BollingerBands) Indicator(candles []*proto.HistoricCandle) (bool, error) {
	convCandles := utils.CandlesToTimeSeries(candles)

	lowerBBI := techan.NewBollingerLowerBandIndicator(techan.NewClosePriceIndicator(convCandles), bb.Window, bb.Sigma)
	upperBBI := techan.NewBollingerUpperBandIndicator(techan.NewClosePriceIndicator(convCandles), bb.Window, bb.Sigma)
	calcLower := lowerBBI.Calculate(len(candles)).Float()
	calcUpper := upperBBI.Calculate(len(candles)).Float()
	closePrice := convCandles.LastCandle().ClosePrice.Float()

	//при пробитии ценой нижней границы - покупка
	if closePrice < calcLower {
		return true, nil
	}

	//при пробитии ценой верхней границы - продажа
	if closePrice > calcUpper {
		return false, nil
	}
	return false, errors.New("unknown state of indicator")
}

func (bb *BollingerBands) GetCandleInterval() int {
	return bb.CandleIntervalInHours
}

func (bb *BollingerBands) GetAnalyzeInterval() int {
	return bb.AnalyzeIntervalInHours
}

func (bb *BollingerBands) GetStrategyParamByString() string {
	return fmt.Sprintf("Bollinger Bands: CandleIntervalHours: %d AnalyzeIntervalHours: %d Window: %d Sigma %f\n", bb.CandleIntervalInHours, bb.AnalyzeIntervalInHours, bb.Window, bb.Sigma)
}
