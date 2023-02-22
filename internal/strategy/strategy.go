package strategy

import (
	"github.com/sdcoffey/techan"
	"io"
	"someshit/pkg/proto"
)

type IndicatorSignal struct {
	Changed bool
	Value   bool
}

type Strategy interface {
	Indicator(candles []*proto.HistoricCandle) (IndicatorSignal, error)
	GetCandleInterval() int
	GetAnalyzeInterval() int
	GetStrategyParamByString() string
	DataPlot(*techan.TimeSeries) (io.WriterTo, error)
}
