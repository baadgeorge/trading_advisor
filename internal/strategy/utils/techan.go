package utils

import (
	"fmt"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
	"someshit/pkg/proto"
)

type CandleStruct struct {
	Period techan.TimePeriod
	Value  float64
}

/*type myTicks struct{
	Ticker plot.Ticker
	Format string
	Time func(t float64) time.Time
}

func (myTicks) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := min; i < max; i++ {
		if math.Mod(i, 5) == 0 {
			ticks = append(ticks, plot.Tick{Value: i, Label: strconv})
		}

	}
}*/

func CandlesToTimeSeries(candles []*proto.HistoricCandle) *techan.TimeSeries {
	var techanCandles []*techan.Candle
	//если последняя свеча не полная, то не добавляем ее
	for i, c := range candles {
		if i == len(candles)-1 && !candles[i].IsComplete {
			break
		}

		tc := &techan.Candle{
			Period: techan.TimePeriod{
				Start: c.Time.AsTime(),
				End:   candles[i+1].Time.AsTime(),
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
