package sdk

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"someshit/pkg/proto"
)

type InstrumentsInterface interface {
	// Метод получения расписания торгов на торговых площадках
	TradingSchedules(exchange string, from, to *timestamp.Timestamp) ([]*proto.TradingSchedule, error)
	// Метод получения облигации по её идентификатору
	BondBy(filters proto.InstrumentRequest) (*proto.Bond, error)
	// Метод получения списка облигаций
	Bonds(status proto.InstrumentStatus) ([]*proto.Bond, error)
	// Метод получения валюты по её идентификатору
	CurrencyBy(filters proto.InstrumentRequest) (*proto.Currency, error)
	// Метод получения списка валют
	Currencies(status proto.InstrumentStatus) ([]*proto.Currency, error)
	// Метод получения инвестиционного фонда по его идентификатору
	EtfBy(filters proto.InstrumentRequest) (*proto.Etf, error)
	// Метод получения списка инвестиционных фондов
	Etfs(status proto.InstrumentStatus) ([]*proto.Etf, error)
	// Метод получения фьючерса по его идентификатору
	FutureBy(filters proto.InstrumentRequest) (*proto.Future, error)
	// Метод получения списка фьючерсов
	Futures(status proto.InstrumentStatus) ([]*proto.Future, error)
	// Метод получения акции по её идентификатору
	ShareBy(filters proto.InstrumentRequest) (*proto.Share, error)
	// Метод получения списка акций
	Shares(status proto.InstrumentStatus) ([]*proto.Share, error)
	// Метод получения основной информации об инструменте
	GetInstrumentBy(filters proto.InstrumentRequest) (*proto.Instrument, error)
	// Метод получения актива по его идентификатору
	GetAssetBy(assetID string) (*proto.AssetFull, error)
	// Метод получения списка активов
	GetAssets() ([]*proto.Asset, error)
	// Метод получения списка избранных инструментов
	GetFavorites() ([]*proto.FavoriteInstrument, error)
}

type InstrumentsService struct {
	client proto.InstrumentsServiceClient
	token  string
}

func NewInstrumentsService(conn *grpc.ClientConn, tkn string) *InstrumentsService {

	client := proto.NewInstrumentsServiceClient(conn)
	return &InstrumentsService{client: client, token: tkn}
}

func (is InstrumentsService) TradingSchedules(exchange string, from, to *timestamp.Timestamp) ([]*proto.TradingSchedule, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()
	res, err := is.client.TradingSchedules(ctx, &proto.TradingSchedulesRequest{
		Exchange: exchange,
		From:     from,
		To:       to,
	})
	if err != nil {
		return nil, err
	}

	return res.Exchanges, nil
}

func (is InstrumentsService) BondBy(filters proto.InstrumentRequest) (*proto.Bond, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.BondBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

// INSTRUMENT_STATUS_BASE
func (is InstrumentsService) Bonds(status proto.InstrumentStatus) ([]*proto.Bond, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.Bonds(ctx, &proto.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) CurrencyBy(filters proto.InstrumentRequest) (*proto.Currency, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.CurrencyBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Currencies(status proto.InstrumentStatus) ([]*proto.Currency, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.Currencies(ctx, &proto.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) EtfBy(filters proto.InstrumentRequest) (*proto.Etf, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.EtfBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Etfs(status proto.InstrumentStatus) ([]*proto.Etf, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.Etfs(ctx, &proto.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) FutureBy(filters proto.InstrumentRequest) (*proto.Future, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.FutureBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Futures(status proto.InstrumentStatus) ([]*proto.Future, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.Futures(ctx, &proto.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) ShareBy(filters proto.InstrumentRequest) (*proto.Share, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.ShareBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) Shares(status proto.InstrumentStatus) ([]*proto.Share, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.Shares(ctx, &proto.InstrumentsRequest{
		InstrumentStatus: status,
	})
	if err != nil {
		return nil, err
	}

	return res.Instruments, nil
}

func (is InstrumentsService) GetInstrumentBy(filters proto.InstrumentRequest) (*proto.Instrument, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.GetInstrumentBy(ctx, &filters)
	if err != nil {
		return nil, err
	}

	return res.Instrument, nil
}

func (is InstrumentsService) GetAssetBy(assetID string) (*proto.AssetFull, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.GetAssetBy(ctx, &proto.AssetRequest{
		Id: assetID,
	})
	if err != nil {
		return nil, err
	}

	return res.Asset, nil
}

func (is InstrumentsService) GetAssets() ([]*proto.Asset, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.GetAssets(ctx, &proto.AssetsRequest{})
	if err != nil {
		return nil, err
	}

	return res.Assets, nil
}

func (is InstrumentsService) GetFavorites() ([]*proto.FavoriteInstrument, error) {
	ctx, cancel := createRequestContext(is.token)
	defer cancel()

	res, err := is.client.GetFavorites(ctx, &proto.GetFavoritesRequest{})
	if err != nil {
		return nil, err
	}

	return res.FavoriteInstruments, nil
}
