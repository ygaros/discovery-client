package cache

import (
	"fmt"
	"log"

	"github.com/ygaros/discovery-client/dto"
)

type sliceCache struct {
	storage []dto.Service
}

func (c *sliceCache) GetServiceInfo(serviceName string) (dto.Service, error) {
	for _, service := range c.storage {
		if service.Name == serviceName {
			return service, nil
		}
	}
	return dto.Service{}, fmt.Errorf("[err]service %s not cached yet", serviceName)
}

func (c *sliceCache) SaveServiceInfo(service dto.Service) error {
	if _, err := c.GetServiceInfo(service.Name); err != nil {
		log.Printf("saving service to cache %v\n", service)
		c.storage = append(c.storage, service)
		return nil
	}
	return fmt.Errorf("[err] service %s cached already", service.Name)
}
func (c *sliceCache) Replace(services []dto.Service) {
	c.storage = services
	log.Printf("updated cache size %d\n", len(c.storage))
}

func NewCache() Cache {
	return &sliceCache{
		storage: make([]dto.Service, 0),
	}
}
