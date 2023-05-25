package util_rpc

import (
	"github.com/tanenking/svrframe/constants"

	"google.golang.org/grpc"
)

var (
	grpcServer *grpc.Server
)

func init() {
	var options = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(constants.RpcSendRecvMaxSize),
		grpc.MaxSendMsgSize(constants.RpcSendRecvMaxSize),
	}

	grpcServer = grpc.NewServer(options...)
}
