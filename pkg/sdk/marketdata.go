package sdk

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"someshit/pkg/proto"
	"time"
)

type MarketDataInterface interface {
	// The method of requesting historical candlesticks by instrument.
	GetCandles(figi string, from, to *timestamp.Timestamp, interval proto.CandleInterval) ([]*proto.HistoricCandle, error)
	// The method of requesting the latest prices for instruments.
	GetLastPrices(figi []string) ([]*proto.LastPrice, error)
	// The method of requesting the status of trading on instruments.
	GetTradingStatus(figi string) (*proto.GetTradingStatusResponse, error)
}

type MarketDataService struct {
	client proto.MarketDataServiceClient
	token  string
}

func NewMarketDataService(conn *grpc.ClientConn, tkn string) *MarketDataService {

	client := proto.NewMarketDataServiceClient(conn)
	return &MarketDataService{client: client, token: tkn}
}

func (mds MarketDataService) GetCandles(figi string, from, to time.Time, interval time.Duration) ([]*proto.HistoricCandle, error) {
	ctx, cancel := createRequestContext(mds.token)
	defer cancel()

	var interval_pb proto.CandleInterval

	switch interval {
	case 24 * time.Hour:
		interval_pb = proto.CandleInterval_CANDLE_INTERVAL_DAY
	case time.Hour:
		interval_pb = proto.CandleInterval_CANDLE_INTERVAL_HOUR
	case time.Minute:
		interval_pb = proto.CandleInterval_CANDLE_INTERVAL_1_MIN
	default:
		return nil, fmt.Errorf("неизвестный интервал свечей: %s\n", interval)
	}

	res, err := mds.client.GetCandles(ctx, &proto.GetCandlesRequest{
		Figi:     figi,
		From:     timestamppb.New(from),
		To:       timestamppb.New(to),
		Interval: interval_pb,
	})
	if err != nil {
		return nil, err
	}

	return res.Candles, nil
}

func (mds MarketDataService) GetLastPrices(figi []string) ([]*proto.LastPrice, error) {
	ctx, cancel := createRequestContext(mds.token)
	defer cancel()

	res, err := mds.client.GetLastPrices(ctx, &proto.GetLastPricesRequest{
		Figi: figi,
	})
	if err != nil {
		return nil, err
	}

	return res.LastPrices, nil
}

func (mds MarketDataService) GetTradingStatus(figi string) (*proto.GetTradingStatusResponse, error) {
	ctx, cancel := createRequestContext(mds.token)
	defer cancel()

	res, err := mds.client.GetTradingStatus(ctx, &proto.GetTradingStatusRequest{
		Figi: figi,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
