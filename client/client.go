package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const appJson string = "application/json"

type clientRequestDTO struct {
	Name   string `json:"name"`
	Port   int    `json:"port"`
	Domain string `json:"domain"`
}
type Client interface {
	HeartBeat() error
	Register() error
}
type client struct {
	serverUrl  string
	serverPort int

	serviceName string
	serviceUrl  string
	servicePort int
}

func (c *client) HeartBeat() error {
	url := fmt.Sprintf("http://%s:%d/heartbeat", c.serverUrl, c.serverPort)
	bodyData := clientRequestDTO{
		Name:   c.serviceName,
		Port:   c.servicePort,
		Domain: c.serviceUrl,
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
		return errors.New("service not registered")
	}
	return nil
}
func (c *client) Register() error {
	url := fmt.Sprintf("http://%s:%d/register", c.serverUrl, c.serverPort)
	bodyData := clientRequestDTO{
		Name:   c.serviceName,
		Port:   c.servicePort,
		Domain: c.serviceUrl,
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
func NewClient(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int) Client {
	return &client{
		serverUrl:  discoveryServerUrl,
		serverPort: discoveryServerPort,

		serviceName: serviceName,
		serviceUrl:  serviceUrl,
		servicePort: servicePort,
	}
}

func NewClientOneTimeRegister(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int) error {
	client := NewClient(discoveryServerUrl, discoveryServerPort, serviceName, serviceUrl, servicePort)
	return client.Register()
}
func NewClientAndHeartBeat(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int) error {
	client := NewClient(discoveryServerUrl, discoveryServerPort, serviceName, serviceUrl, servicePort)
	registerLoop(client)
	go doHeartBeat(client)
	return nil
}

func registerLoop(client Client) {
	for range time.Tick(5 * time.Second) {
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
