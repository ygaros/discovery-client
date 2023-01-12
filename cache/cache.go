package cache

import (
	"fmt"
	"ygaros-discovery-client/dto"
)

type Cache interface {
	GetServiceInfo(serviceName string) (dto.ServiceDTO, error)
	SaveServiceInfo(service dto.ServiceDTO) error
	Replace(services []dto.ServiceDTO)
}
type cache struct {
	storage []dto.ServiceDTO
}

func (c *cache) GetServiceInfo(serviceName string) (dto.ServiceDTO, error) {
	for _, service := range c.storage {
		if service.Name == serviceName {
			return service, nil
		}
	}
	return dto.ServiceDTO{}, fmt.Errorf("[err]service %s not cached yet", serviceName)
}

func (c *cache) SaveServiceInfo(service dto.ServiceDTO) error {
	if _, err := c.GetServiceInfo(service.Name); err != nil {
		c.storage = append(c.storage, service)
		return nil
	}
	return fmt.Errorf("[err] service %s cached already", service.Name)
}
func (c *cache) Replace(services []dto.ServiceDTO) {
	c.storage = services
}

func NewCache() Cache {
	return &cache{}
}
