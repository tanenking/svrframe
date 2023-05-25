package config

var serviceConfig *ServiceConfig

func init() {
	serviceConfig = &ServiceConfig{
		Service: new(ServiceInfo),
	}
}
