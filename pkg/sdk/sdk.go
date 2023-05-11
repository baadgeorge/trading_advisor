package sdk

type ServicePool struct {
	InstrumentsService InstrumentsInterface
	MarketDataService  MarketDataInterface
	UsersService       UsersServiceInterface
}

// метод проверки токена на валидность
func (sp *ServicePool) TokenIsValid() bool {
	_, err := sp.UsersService.GetInfo()
	if err != nil {
		return false
	}
	return true
}

// функция создания экземпляра сервисов с указанным токеном
func NewServicePool(token string) (*ServicePool, error) {
	conn, err := clientConnection()
	if err != nil {
		return nil, err
	}

	return &ServicePool{
		InstrumentsService: *NewInstrumentsService(conn, token),
		MarketDataService:  *NewMarketDataService(conn, token),
		UsersService:       *NewUsersService(conn, token),
	}, nil
}
