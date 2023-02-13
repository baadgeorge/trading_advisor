package sdk

// non-stream services
type ServicePool struct {
	InstrumentsService InstrumentsService
	MarketDataService  MarketDataService
	UsersService       UsersService
}

func NewServicePool(token string) *ServicePool {
	return &ServicePool{
		InstrumentsService: *NewInstrumentsService(token),
		MarketDataService:  *NewMarketDataService(token),
		UsersService:       *NewUsersService(token),
	}
}
