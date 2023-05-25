package application

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime/trace"
	"syscall"
	"time"

	"github.com/tanenking/svrframe/config"
	"github.com/tanenking/svrframe/constants"
	"github.com/tanenking/svrframe/helper"
	"github.com/tanenking/svrframe/logx"
	"github.com/tanenking/svrframe/mysqlx"
	"github.com/tanenking/svrframe/redisx"
)

// deify
var logo = `


  ██████╗ ███████╗██╗███████╗██╗   ██╗
  ██╔══██╗██╔════╝██║██╔════╝╚██╗ ██╔╝
  ██║  ██║█████╗  ██║█████╗   ╚████╔╝ 
  ██║  ██║██╔══╝  ██║██╔══╝    ╚██╔╝  
  ██████╔╝███████╗██║██║        ██║   
  ╚═════╝ ╚══════╝╚═╝╚═╝        ╚═╝   
`
var topLine = `┌───────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└───────────────────────────────────────────────────┘`

func InitProgram(service_type string, etcdNameFunc func() string, customConfigFunc func() error) bool {

	// if constants.IsCoredump() {
	// 	debug.SetTraceback("crash")
	// }
	// if doFork() {
	// 	return false
	// }

	constants.Service_Type = service_type

	if err := config.ParseConfig(etcdNameFunc, customConfigFunc); err != nil {
		return false
	}

	logx.InitLogx()
	rand.Seed(time.Now().Unix())

	_logo := logo + "\n"
	_logo += topLine + "\n"
	_logo += fmt.Sprintf("%s [coder] martin                                    %s", borderLine, borderLine) + "\n"
	_logo += fmt.Sprintf("%s [time] 2022-01-29                                 %s", borderLine, borderLine) + "\n"
	_logo += bottomLine + "\n"

	logx.InfoF("%s", _logo)

	mysqlx.InitMysqlHelper()
	redisx.InitRedisHelper()

	return true
}

func ProgramRunning() {
	defer fmt.Println("server close")
	defer func() {
		constants.GetServiceStopListener().NotifyAllListeners()
		constants.GetServiceStopWaitGroup().Wait()
		constants.GetServiceStopListener().Clear()
		//stopTrace()
	}()
	//startTrace()

	logx.InfoF("服务启动成功")
	if !writePid() {
		return
	}
	//wait for exit
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig, ok := <-ch:
			if ok {
				logx.InfoF("signal receive: %v", sig)
				switch sig {
				case syscall.SIGINT:
					return
				case syscall.SIGTERM:
					return
				}
			}
		default:
			helper.GetGlobalTimer().Update()
		}

		time.Sleep(time.Millisecond * 10)
	}
}
func EnablePprof(port int) {
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}
func startTrace() {
	if !constants.IsDebug() {
		return
	}
	file, _ := os.Create(logx.GetTraceFile())
	trace.Start(file)
}
func stopTrace() {
	if !constants.IsDebug() {
		return
	}
	trace.Stop()
}
