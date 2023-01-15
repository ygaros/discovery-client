package cache

import (
	"github.com/ygaros/discovery-client/dto"
)

type Cache interface {
	GetServiceInfo(serviceName string) (dto.Service, error)
	SaveServiceInfo(service dto.Service) error
	Replace(services []dto.Service)
}
