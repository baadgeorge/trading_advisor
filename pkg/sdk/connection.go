package sdk

import (
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const ApiURL = "invest-public-api.tinkoff.ru:443"

func clientConnection() (*grpc.ClientConn, error) {
	tlsConfig := tls.Config{InsecureSkipVerify: true}
	return grpc.Dial(ApiURL, grpc.WithTransportCredentials(credentials.NewTLS(&tlsConfig)))
}
