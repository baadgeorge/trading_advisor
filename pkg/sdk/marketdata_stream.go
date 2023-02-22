package sdk

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"someshit/pkg/proto"
)

type MarketDataStreamInterface interface {
	// Recv listens for incoming messages and block until first one is received.
	Recv() (*proto.MarketDataResponse, error)
	// Send puts proto.MarketDataRequest into a stream.
	Send(request *proto.MarketDataRequest) error
}

type MarketDataStream struct {
	client proto.MarketDataStreamServiceClient
	stream proto.MarketDataStreamService_MarketDataStreamClient
}

func NewMarketDataStream(conn *grpc.ClientConn, token string) *MarketDataStream {
	/*conn, err := clientConnection()
	if err != nil {
		logrus.Fatal(err.Error())
	}*/

	client := proto.NewMarketDataStreamServiceClient(conn)
	ctx := createStreamContext(token)

	stream, err := client.MarketDataStream(ctx)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	return &MarketDataStream{client: client, stream: stream}
}

func (mds MarketDataStream) Recv() (*proto.MarketDataResponse, error) {
	return mds.stream.Recv()
}

func (mds MarketDataStream) Send(request *proto.MarketDataRequest) error {
	return mds.stream.Send(request)
}
