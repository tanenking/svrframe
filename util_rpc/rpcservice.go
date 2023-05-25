package util_rpc

import (
	"fmt"
	"net"

	config "github.com/tanenking/svrframe/config"
	"github.com/tanenking/svrframe/constants"
	"github.com/tanenking/svrframe/logx"

	"google.golang.org/grpc"
)

func update() {
	go func() {
		exit_ch := constants.GetServiceStopListener().AddListener()
		<-exit_ch.Done()
		grpcServer.Stop()
		grpcServer = nil
	}()
}

func StartRpcService(regfunc func(server *grpc.Server)) error {
	serviceConfig := config.GetServiceConfig()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", serviceConfig.Service.RpcPort))
	if err != nil {
		logx.ErrorF("failed to listen: %v", err)
		return err
	}
	if serviceConfig.Service.RpcPort == 0 {
		addr, err := net.ResolveTCPAddr(lis.Addr().Network(), lis.Addr().String())
		if err != nil {
			logx.ErrorF("%v", err)
			return err
		}
		serviceConfig.Service.RpcPort = uint16(addr.Port)
	}
	logx.InfoF("RPC SERVER [ %s ] RUNNING", constants.Service_Type)

	regfunc(grpcServer)
	go func() {
		grpcServer.Serve(lis)
	}()
	update()

	return nil
}
