package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"ygaros-discovery-client/cache"
	"ygaros-discovery-client/dto"
)

const (
	appJson = "application/json"
	HTTPS   = "https://"
	HTTP    = "http://"
)

type registerDTO struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Secure bool   `json:"secure"`
}

type Client interface {
	HeartBeat() error
	Register() error
	GetService(serviceName string) (*dto.ServiceDTO, error)
	StartFetchingServices(duration time.Duration)
}

type client struct {
	serverUrl string

	serviceName string
	serviceUrl  string
	servicePort int
	secure      bool
	cache       cache.Cache
}

func (c *client) HeartBeat() error {
	url := fmt.Sprintf("%s/heartbeat", c.serverUrl)
	bodyData := registerDTO{
		Name:   c.serviceName,
		Url:    c.serviceUrl,
		Secure: c.secure,
	}
	log.Println("HeartBeating on discovery server..")
	marshalledBody, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}
	post, err := http.Post(
		url,
		appJson,
		bytes.NewBuffer(marshalledBody),
	)
	if err != nil {
		return err
	}
	if post.StatusCode != http.StatusOK {
		return errors.New("[err] service not registered")
	}
	return nil
}

func (c *client) Register() error {
	url := fmt.Sprintf("%s/register", c.serverUrl)
	bodyData := registerDTO{
		Name:   c.serviceName,
		Url:    c.serviceUrl,
		Secure: c.secure,
	}
	log.Println("Registering on discovery server...")
	marshalledBody, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}
	_, err = http.Post(
		url,
		appJson,
		bytes.NewBuffer(marshalledBody),
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) GetService(serviceName string) (*dto.ServiceDTO, error) {
	if cached, err := c.cache.GetServiceInfo(serviceName); err == nil {
		log.Printf("Getting service url from cache %s\n", cached.Name)
		return &cached, err
	}
	url := fmt.Sprintf("%s/service?serviceName=%s", c.serverUrl, serviceName)
	response, _ := http.Get(url)
	service := dto.ServiceDTO{}
	if response.StatusCode != http.StatusOK {
		return &service, fmt.Errorf("[err] service %s isnt registered in discovery server", serviceName)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return &service, err
	}
	defer response.Body.Close()
	err = json.Unmarshal(body, &service)
	if err != nil {
		return &service, err
	}
	log.Printf("Saving discovered service to cache %s\n", service.Name)
	c.cache.SaveServiceInfo(service)
	return &service, nil
}

func (c *client) StartFetchingServices(duration time.Duration) {
	ticker := time.NewTicker(duration)
	for range ticker.C {
		services, err := c.getListOfServices()
		if err != nil {
			for _, service := range services {
				c.cache.SaveServiceInfo(service)
			}
		}
	}
}

func (c *client) getListOfServices() ([]dto.ServiceDTO, error) {
	url := fmt.Sprintf("%s/list", c.serverUrl)
	response, err := http.Get(url)
	var service []dto.ServiceDTO
	if err != nil {
		return service, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return service, err
	}
	defer response.Body.Close()
	err = json.Unmarshal(body, &service)
	if err != nil {
		return service, err
	}
	return service, nil
}

func NewClientOneTimeRegister(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int, secure bool) error {
	client := NewClient(
		discoveryServerUrl, discoveryServerPort,
		serviceName, serviceUrl, servicePort, secure,
	)
	return client.Register()
}

func NewClientAndHeartBeat(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int, secure bool) (Client, error) {
	client := NewClient(
		discoveryServerUrl, discoveryServerPort,
		serviceName, serviceUrl, servicePort, secure,
	)
	registerLoop(client)
	go doHeartBeat(client)
	go client.StartFetchingServices(30 * time.Second)
	return client, nil
}
func NewClient(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int, secure bool) Client {
	return &client{
		serverUrl: fmt.Sprintf("http://%s:%d", discoveryServerUrl, discoveryServerPort),

		serviceName: serviceName,
		serviceUrl:  fmt.Sprintf("%s:%d", serviceUrl, servicePort),
		servicePort: servicePort,
		cache:       cache.NewCache(),
	}
}

func registerLoop(client Client) {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		err := client.Register()
		if err == nil {
			log.Println("Registered successfully in discovery server!")
			break
		} else {
			log.Println(err)
		}
	}
}

func doHeartBeat(client Client) {
	log.Println("Heartbeat started at 1/1min rate")
	for range time.Tick(time.Minute) {
		err := client.HeartBeat()
		if err != nil {
			log.Println(err)
			registerLoop(client)
		}
	}
}
