package strategy

import (
	"final/pkg/proto"
	"github.com/sdcoffey/techan"
)

type IndicatorSignal struct {
	Changed bool
	Value   bool
}

type Strategy interface {
	//метод вычисления индикатора
	Indicator(candles []*proto.HistoricCandle) (IndicatorSignal, error)
	//метод получения интервала свечей
	GetCandleInterval() int
	//метод получения размера окна
	GetAnalyzeInterval() int
	//метод получения параметров стратегии в виде строки
	GetStrategyParamByString() string
	//метод получения графического представления свечй и индикатора
	//в виде байтов
	DataPlot(*techan.TimeSeries) ([][]byte, error)
}
