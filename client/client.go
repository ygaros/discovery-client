package client

import (
	"time"

	"github.com/ygaros/discovery-client/dto"
)

type Client interface {
	HeartBeat() error
	Register() error
	GetService(serviceName string) (*dto.Service, error)
	StartFetchingServices(duration time.Duration)
	GetListOfServices() ([]dto.Service, error)
}
