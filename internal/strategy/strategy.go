package strategy

import (
	"github.com/sdcoffey/techan"
	"io"
	"someshit/pkg/proto"
)

type Strategy interface {
	Indicator(candles []*proto.HistoricCandle) (bool, error)
	GetCandleInterval() int
	GetAnalyzeInterval() int
	GetStrategyParamByString() string
	DataPlot(*techan.TimeSeries) (io.WriterTo, error)
}
