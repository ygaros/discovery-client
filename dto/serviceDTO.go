package dto

import (
	"time"
)

type ServiceDTO struct {
	Name               string    `json:"name"`
	Url                string    `json:"url"`
	LastHeartBeatCheck time.Time `json:"lastHeartBeatCheck"`
}
