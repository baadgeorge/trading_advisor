package utils

import (
	"final/pkg/proto"
	"fmt"

	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

// функция конвертации типа свечей для выичсления значения индикатора
func CandlesToTimeSeries(candles []*proto.HistoricCandle) *techan.TimeSeries {
	var techanCandles []*techan.Candle
	//если последняя свеча не полная, то не анализируем ее
	for i, c := range candles {
		if i == len(candles)-1 && !candles[i].IsComplete {
			break
		}
		if i == 0 {
			continue
		}
		tc := &techan.Candle{
			Period: techan.TimePeriod{
				Start: candles[i-1].Time.AsTime().Local(),
				End:   c.Time.AsTime().Local(),
			},
			Volume:     big.NewFromInt(int(c.Volume)),
			OpenPrice:  big.NewFromString(fmt.Sprintf("%d.%d", c.Close.Units, abs(c.Close.Nano))),
			ClosePrice: big.NewFromString(fmt.Sprintf("%d.%d", c.Close.Units, abs(c.Close.Nano))),
			MaxPrice:   big.NewFromString(fmt.Sprintf("%d.%d", c.High.Units, abs(c.High.Nano))),
			MinPrice:   big.NewFromString(fmt.Sprintf("%d.%d", c.Low.Units, abs(c.Low.Nano))),
		}
		techanCandles = append(techanCandles, tc)
	}
	return &techan.TimeSeries{Candles: techanCandles}
}

func abs(x int32) int32 {
	return absDiff(x, 0)
}

func absDiff(x, y int32) int32 {
	if x < y {
		return y - x
	}
	return x - y
}
