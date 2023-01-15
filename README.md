# Ygaros Discovery Client

Basic discovery client for microservice architecture with in memory (slice) caching

[**SERVER**](https://github.com/ygaros/discovery-server)

*The idea of this project is to try mimic the `Eureka-Client` from `Spring-Cloud` in **Go**.*

### To get started

```
go get github.com/ygaros/discovery-client
```

*Fully functional **gRPC** discovery-client with in memory caching.*


```
    discoveryServerUrl := "localhost"
	dicoveryServerPort := 7654

	serviceName := "service-on-8000"
	serviceUrl := "localhost"
	servicePort := 8000
	isHttps := false

	c, err = client.NewClient(
		discoveryServerUrl,
		dicoveryServerPort,
		serviceName,
		serviceUrl,
		servicePort,
		isHttps,
	)
```
*Client register to discovery-server `via` data provided to `NewClient()` and **HeartBeats** every 30 seconds.*

To get data about another service registered in discovery server.

```
    service, err := c.GetService("another-registered-service-name")
```
*`GetService()` returns `*dto.Service`.*

```
type Service struct {
	Name               string   
	Url                string   
	LastHeartBeatCheck time.Time 
}
```
*Url field id ready-to-call-url with correct protocol (HTTP/HTTPS) inluded.*

### Default timers

*Every 30 seconds client HeartBeats to discovery-server.*\
*Every 30 seconds client is calling discovery-server to get all registered services for caching.*