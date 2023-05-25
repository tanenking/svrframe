package constants

import (
	"context"
	"net"
	"os"
	"runtime"
	"strings"

	"google.golang.org/grpc/peer"
)

type (
	ctx_key string
)

const (
	RpcSendRecvMaxSize = 1024 * 1024 * 64
)

const (
	RequestParam ctx_key = "requestParam"
)

const (
	ServiceMode_TEST   = "TEST"
	ServiceMode_FORMAL = "FORMAL"
)

const (
	RuntimeMode_Debug   = "debug"
	RuntimeMode_Release = "release"
)

const (
	Env_ServiceHost = "service_host"
	Env_Gopath      = "GOPATH"
	Env_RuntimeMode = "runtime_mode"
	Env_ServiceMode = "service_mode"
	Env_Coredump    = "coredump"
	Env_LogLevel    = "log_level"
	Env_LogPath     = "log_path"
	Env_LogRuntime  = "log_runtime"
	Env_CfgRootPath = "cfg_path"
)

const (
	TimeFormatString   = "2006-01-02 15:04:05"
	TimeFormat20060102 = "20060102"
)

const (
	TenThousandthRatio = 0.0001
)

// /////////////////////////////////////////////////////////////////////////////
func GetSystem() string {
	return runtime.GOOS
}
func IsWindowsSystem() bool {
	return GetSystem() == "windows"
}

func GetServiceMode() string {
	mode := strings.ToUpper(os.Getenv(Env_ServiceMode))
	if len(mode) <= 0 || !IsValidServiceMode(mode) {
		mode = ServiceMode_TEST
	}
	return mode
}

func IsValidServiceMode(mode string) bool {
	_, ok := Service_mode_map[mode]
	return ok
}

func IsCoredump() bool {
	coredump_mode := os.Getenv(Env_Coredump)
	return len(coredump_mode) > 0
}

func IsDebug() bool {
	runtime_mode := os.Getenv(Env_RuntimeMode)
	return strings.ToLower(runtime_mode) == RuntimeMode_Debug
}

func GetServiceHost() string {
	addr := os.Getenv(Env_ServiceHost)
	if len(addr) > 0 && net.ParseIP(addr) != nil {
		return addr
	}

	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	return "127.0.0.1"
}

func GetPeerAddr(ctx context.Context) string {
	var addr string
	if pr, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			addr = tcpAddr.IP.String()
		} else {
			addr = pr.Addr.String()
		}
	}
	return addr
}
