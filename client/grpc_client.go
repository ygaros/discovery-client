package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ygaros/discovery-client/cache"
	"github.com/ygaros/discovery-client/dto"
	proto "github.com/ygaros/discovery-client/gen/proto"

	"google.golang.org/grpc"
)

type grpcClient struct {
	serverUrl string

	serviceName string
	serviceUrl  string
	servicePort int
	secure      bool
	cache       cache.Cache
	ctx         context.Context
	client      proto.DiscoveryClient
}

func (gc *grpcClient) HeartBeat() error {
	service := proto.Service{
		Name:   gc.serviceName,
		Url:    gc.serviceUrl,
		Secure: gc.secure,
	}
	log.Printf("HeartBeating %s to discovery server\n", service.Name)
	_, err := gc.client.HeartBeat(gc.ctx, &service)
	return err
}

func (gc *grpcClient) Register() error {
	service := proto.Service{
		Name:   gc.serviceName,
		Url:    gc.serviceUrl,
		Secure: gc.secure,
	}
	log.Printf("Registering %s\n", service.Name)
	_, err := gc.client.AddService(context.Background(), &service)
	return err
}

func (gc *grpcClient) GetService(serviceName string) (*dto.Service, error) {
	if cached, err := gc.cache.GetServiceInfo(serviceName); err == nil {
		log.Printf("Getting Service data for %s from cache\n", serviceName)
		return &cached, err
	}
	request := proto.GetServiceRequest{
		ServiceName: serviceName,
	}
	log.Printf("Getting Service data for %s from discovery server\n", serviceName)
	serviceResponse, err := gc.client.GetService(gc.ctx, &request)
	if err != nil {
		return &dto.Service{}, err
	}
	service, err := dto.ToService(serviceResponse)
	return &service, err
}

func (gc *grpcClient) StartFetchingServices(duration time.Duration) {
	log.Println("Starting fetching all services goroutine")
	ticker := time.NewTicker(duration)
	for range ticker.C {
		services, err := gc.GetListOfServices()
		if err == nil {
			gc.cache.Replace(services)
		}
	}
}

func (gc *grpcClient) GetListOfServices() (services []dto.Service, err error) {
	response, err := gc.client.ListServices(gc.ctx, &proto.Empty{})
	if err != nil {
		return nil, err
	}
	for _, service := range response.Services {
		if parsed, err := dto.ToService(service); err == nil {
			services = append(services, parsed)
		}
	}
	return services, err
}

func NewGrpcClient(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int, secure bool) Client {
	serverUrl := fmt.Sprintf("%s:%d", discoveryServerUrl, discoveryServerPort)
	conn, err := grpc.Dial(serverUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	return &grpcClient{
		serverUrl: serverUrl,

		serviceName: serviceName,
		serviceUrl:  fmt.Sprintf("%s:%d", serviceUrl, servicePort),
		servicePort: servicePort,
		client:      proto.NewDiscoveryClient(conn),
		ctx:         context.Background(),
		cache:       cache.NewCache(),
	}
}
func NewGrpcClientOneTimeRegister(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int, secure bool) error {
	client := NewGrpcClient(
		discoveryServerUrl, discoveryServerPort,
		serviceName, serviceUrl, servicePort, secure,
	)
	return client.Register()
}
func NewGrpcClientAndHeartBeat(discoveryServerUrl string, discoveryServerPort int,
	serviceName, serviceUrl string, servicePort int, secure bool) (Client, error) {
	client := NewGrpcClient(
		discoveryServerUrl, discoveryServerPort,
		serviceName, serviceUrl, servicePort, secure,
	)
	RegisterLoop(client)
	go DoHeartBeat(client)
	go client.StartFetchingServices(30 * time.Second)
	return client, nil
}
