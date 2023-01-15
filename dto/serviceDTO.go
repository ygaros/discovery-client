package dto

import (
	"time"

	proto "github.com/ygaros/discovery-client/gen/proto"
)

const TIME_FORMAT = "2006-01-02 15:04:05.999999999 -0700 MST"

type Service struct {
	Name               string    `json:"name"`
	Url                string    `json:"url"`
	LastHeartBeatCheck time.Time `json:"lastHeartBeatCheck"`
}
type Register struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Secure bool   `json:"secure"`
}

func ToService(service *proto.ServiceWithHeartBeat) (Service, error) {
	time, err := time.Parse(TIME_FORMAT, service.LastHeartBeat)
	if err != nil {
		return Service{}, err
	}
	return Service{
		Name:               service.Name,
		Url:                service.Url,
		LastHeartBeatCheck: time,
	}, nil
}
