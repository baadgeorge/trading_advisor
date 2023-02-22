package sdk

type ServicePool struct {
	InstrumentsService InstrumentsService
	MarketDataService  MarketDataService
	MarketDataStream   MarketDataStream
	UsersService       UsersService
}

func NewServicePool(token string) (*ServicePool, error) {
	conn, err := clientConnection()
	if err != nil {
		return nil, err
	}

	return &ServicePool{
		InstrumentsService: *NewInstrumentsService(conn, token),
		MarketDataService:  *NewMarketDataService(conn, token),
		MarketDataStream:   *NewMarketDataStream(conn, token),
		UsersService:       *NewUsersService(conn, token),
	}, nil
}

/*func (sp ServicePool) ListenNewCandles() {

	mds := NewMarketDataStream()
	recv := proto.SubscribeCandlesRequest{
		SubscriptionAction: 0,
		Instruments:        proto.CandleInstrument{Interval: },
		WaitingClose:       false,
	}
	payload := proto.MarketDataRequest_SubscribeCandlesRequest{SubscribeCandlesRequest: recv}

	for {
		newEvent, err := sp.MarketDataStream.client.
		if err != nil {
			logrus.Errorf("Failed while receiving new candles: %v", err)
			return
		}
		if err != io.EOF {
			logrus.Errorf("Candle stream closed: %v", err)
			return
		}
		if newEvent != nil && newEvent.GetCandle() != nil {
			figi := newEvent.GetCandle().Figi
			interval := newEvent.Ge
		}
	}

}*/
