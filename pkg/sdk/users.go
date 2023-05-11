package sdk

import (
	"final/pkg/proto"
	"google.golang.org/grpc"
)

type UsersServiceInterface interface {
	// Метод получения информации о пользователе
	GetInfo() (*proto.GetInfoResponse, error)
}

type UsersService struct {
	client proto.UsersServiceClient
	token  string
}

func NewUsersService(conn *grpc.ClientConn, tkn string) *UsersService {
	client := proto.NewUsersServiceClient(conn)
	return &UsersService{
		client: client,
		token:  tkn}
}

func (us UsersService) GetInfo() (*proto.GetInfoResponse, error) {
	ctx, cancel := createRequestContext(us.token)
	defer cancel()

	res, err := us.client.GetInfo(ctx, &proto.GetInfoRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
