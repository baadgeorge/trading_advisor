package sdk

import (
	"google.golang.org/grpc"
	"someshit/pkg/proto"
)

type UsersServiceClient interface {
	// The method of receiving user accounts.
	GetAccounts() ([]*proto.Account, error)
	// The method of obtaining information about the user.
	GetInfo() (*proto.GetInfoResponse, error)
}

type UsersService struct {
	client proto.UsersServiceClient
	token  string
}

func NewUsersService(conn *grpc.ClientConn, tkn string) *UsersService {
	//conn, err := clientConnection()

	client := proto.NewUsersServiceClient(conn)
	return &UsersService{
		client: client,
		token:  tkn}
}

func (us UsersService) GetAccounts() ([]*proto.Account, error) {
	ctx, cancel := createRequestContext(us.token)
	defer cancel()

	res, err := us.client.GetAccounts(ctx, &proto.GetAccountsRequest{})
	if err != nil {
		return nil, err
	}

	return res.Accounts, nil
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
