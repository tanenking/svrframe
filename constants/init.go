package constants

import "sync"

var (
	serviceStopListener  *IListenerManager
	serviceStopWaitGroup sync.WaitGroup
	Service_Type         string
	Service_mode_map     map[string]string

	ProjectName string //项目名称
)

func init() {
	serviceStopListener = NewListenerManager()
	serviceStopWaitGroup = sync.WaitGroup{}
	Service_Type = "unknow"

	Service_mode_map = map[string]string{
		ServiceMode_TEST:   ServiceMode_TEST,
		ServiceMode_FORMAL: ServiceMode_FORMAL,
	}
}

func GetServiceStopListener() *IListenerManager {
	return serviceStopListener
}
func GetServiceStopWaitGroup() *sync.WaitGroup {
	return &serviceStopWaitGroup
}
