package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tanenking/svrframe/constants"
)

type ServiceConfig struct {
	ProjectName  string
	Service      *ServiceInfo
	Etcd         *EtcdConfig
	Mysql        map[string]*MysqlConfig
	RedisCluster *RedisClusterConfig
	CustomConfig interface{} //自定义配置内容
}
type EtcdConfig struct {
	EtcdName  string
	Endpoints []string
	WatchList []string
}
type ServiceInfo struct {
	//服务rpc信息
	RpcConfig `json:"Rpc,omitempty"`
	//服务http信息
	HttpConfig `json:"Http,omitempty"`
	//服务http信息
	TcpConfig `json:"Tcp,omitempty"`
	//服务大区ID
	ZoneID uint16
	//服务ID
	ServiceID uint16
	//服务地址
	ServiceHost string
}
type RpcConfig struct {
	RpcPort uint16
}
type TcpConfig struct {
	TcpPort uint16
	MaxConn uint16
	IsWS    bool
}
type HttpConfig struct {
	HttpPort uint16
}
type MysqlConfig struct {
	Name     string
	UserName string
	Password string
	Host     string
	Port     uint16
	Database string
	Charset  string
}
type RedisClusterConfig struct {
	Password string
	Redis    []*RedisConfig
}
type RedisConfig struct {
	Host string
	Port uint16
}

func ParseConfig(etcdNameFunc func() string, customConfigFunc func() error) (err error) {
	cfgFile := flag.String("config", "", "config file path")
	//解析配置文件
	flag.Parse()

	if cfgFile == nil || len(*cfgFile) <= 0 {
		if constants.IsDebug() {
			*cfgFile = "configs/front.json"
		} else {
			fmt.Printf("need cmdline param -config=config file path\n")
			return
		}
	}

	err = readConfig(*cfgFile, etcdNameFunc, customConfigFunc)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	return
}

func readConfig(config_file string, etcdNameFunc func() string, customConfigFunc func() error) (err error) {
	dir, _ := os.Getwd()
	dir = strings.ReplaceAll(dir, "\\", "/")
	fullPath := fmt.Sprintf("%s/%s", dir, config_file)

	filePtr, err := os.Open(fullPath)
	if err != nil {
		paths := strings.Split(dir, "/")
		paths = paths[0 : len(paths)-1]
		dir = strings.Join(paths, "/")
		fmt.Printf("not found config %s, try found in gopath\n", fullPath)
		// dir = os.Getenv(constants.Env_Gopath)
		fullPath = fmt.Sprintf("%s/%s", dir, config_file)
		filePtr, err = os.Open(fullPath)
		if err != nil {
			fmt.Printf("Open file failed [Err:%s]\n", err.Error())
			return
		}
	}
	defer filePtr.Close()

	allbytes, err := io.ReadAll(filePtr)
	if err != nil {
		fmt.Printf("Read file failed [Err:%s]\n", err.Error())
		return
	}

	err = json.Unmarshal(allbytes, serviceConfig)
	if err != nil {
		fmt.Printf("Parse file failed [Err:%s]\n", err.Error())
		return
	}
	if len(serviceConfig.Service.ServiceHost) <= 0 {
		serviceConfig.Service.ServiceHost = constants.GetServiceHost()
	}
	constants.ProjectName = serviceConfig.ProjectName

	if customConfigFunc != nil {
		if err = customConfigFunc(); err != nil {
			return
		}
	}

	watchList := []string{}
	if len(serviceConfig.Etcd.WatchList) > 0 {
		for _, ws := range serviceConfig.Etcd.WatchList {
			watchList = append(watchList, fmt.Sprintf("/%s/%s", constants.ProjectName, ws))
		}
	}
	serviceConfig.Etcd.WatchList = watchList
	if etcdNameFunc != nil {
		serviceConfig.Etcd.EtcdName = etcdNameFunc()
	} else {
		serviceConfig.Etcd.EtcdName = GetEtcdNameDefault()
	}

	if serviceConfig.Mysql != nil && len(serviceConfig.Mysql) > 0 {
		for name, m := range serviceConfig.Mysql {
			if m.Port == 0 || len(m.Database) <= 0 || len(m.UserName) <= 0 || len(m.Password) <= 0 {
				return fmt.Errorf("mysql config has wrong")
			}
			m.Charset = "utf8mb4"
			m.Name = name
		}
	}

	fmt.Printf("Parse file [%s] success \n", config_file)
	return
}
func GetEtcdNameDefault() string {
	return fmt.Sprintf("/%s/%s/z%d/svr-%d", constants.ProjectName, constants.Service_Type, serviceConfig.Service.ZoneID, serviceConfig.Service.ServiceID)
}

func GetServiceConfig() *ServiceConfig {
	return serviceConfig
}
func GetServiceInfo() *ServiceInfo {
	return serviceConfig.Service
}
func GetMysqlConfigs() map[string]*MysqlConfig {
	return serviceConfig.Mysql
}
func GetMysqlConfig(dbName string) *MysqlConfig {
	dbcfg, ok := serviceConfig.Mysql[dbName]
	if !ok {
		return nil
	}
	return dbcfg
}
func GetRedisClusterConfig() *RedisClusterConfig {
	return serviceConfig.RedisCluster
}
func GetEtcdInfo() *EtcdConfig {
	return serviceConfig.Etcd
}
func GetCustomConfig() interface{} {
	return serviceConfig.CustomConfig
}
