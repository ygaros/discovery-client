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

	"github.com/ygaros/discovery-client/cache"
	"github.com/ygaros/discovery-client/dto"
)

const (
	appJson = "application/json"
	HTTPS   = "https://"
	HTTP    = "http://"
)

type httpClient struct {
	serverUrl string

	serviceName string
	serviceUrl  string
	servicePort int
	secure      bool
	cache       cache.Cache
}

func (c *httpClient) HeartBeat() error {
	url := fmt.Sprintf("%s/heartbeat", c.serverUrl)
	bodyData := dto.Register{
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

func (c *httpClient) Register() error {
	url := fmt.Sprintf("%s/register", c.serverUrl)
	bodyData := dto.Register{
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

func (c *httpClient) GetService(serviceName string) (*dto.Service, error) {
	if cached, err := c.cache.GetServiceInfo(serviceName); err == nil {
		log.Printf("Getting service url from cache %s\n", cached.Name)
		return &cached, err
	}
	url := fmt.Sprintf("%s/service?serviceName=%s", c.serverUrl, serviceName)
	response, _ := http.Get(url)
	service := dto.Service{}
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

func (c *httpClient) StartFetchingServices(duration time.Duration) {
	ticker := time.NewTicker(duration)
	for range ticker.C {
		services, err := c.GetListOfServices()
		if err != nil {
			for _, service := range services {
				c.cache.SaveServiceInfo(service)
			}
		}
	}
}

func (c *httpClient) GetListOfServices() ([]dto.Service, error) {
	url := fmt.Sprintf("%s/list", c.serverUrl)
	response, err := http.Get(url)
	var service []dto.Service
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

func NewHttpClientOneTimeRegister(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int, secure bool) error {
	client := NewHttpClient(
		discoveryServerUrl, discoveryServerPort,
		serviceName, serviceUrl, servicePort, secure,
	)
	return client.Register()
}

func NewHttpClientAndHeartBeat(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int, secure bool) (Client, error) {
	client := NewHttpClient(
		discoveryServerUrl, discoveryServerPort,
		serviceName, serviceUrl, servicePort, secure,
	)
	RegisterLoop(client)
	go DoHeartBeat(client)
	go client.StartFetchingServices(30 * time.Second)
	return client, nil
}
func NewHttpClient(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int, secure bool) Client {
	return &httpClient{
		serverUrl: fmt.Sprintf("http://%s:%d", discoveryServerUrl, discoveryServerPort),

		serviceName: serviceName,
		serviceUrl:  fmt.Sprintf("%s:%d", serviceUrl, servicePort),
		servicePort: servicePort,
		cache:       cache.NewCache(),
	}
}
