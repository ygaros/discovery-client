package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ygaros/discovery-client/client"
)

var serverClient client.Client

func main() {
	var err error

	discoveryServerUrl := "localhost"
	dicoveryServerPort := 7654

	serviceName := "service-on-8000"
	serviceUrl := "localhost"
	servicePort := 8000
	isHttps := false

	serverClient, err = client.NewClient(
		discoveryServerUrl,
		dicoveryServerPort,
		serviceName,
		serviceUrl,
		servicePort,
		isHttps,
	)
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", GetMapping)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", servicePort), nil))

}

type Response struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

func GetMapping(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Time:    time.Now(),
		Message: "Merged message:",
	}
	serviceArray := [3]string{
		"service-on-8001",
		"service-on-8002",
		"service-on-8003",
	}
	for _, serviceName := range serviceArray {
		serviceData, err := serverClient.GetService(serviceName)
		if err != nil {
			log.Println(err)
			continue
		}
		respo, err := DoGetRequest(serviceData.Url)
		if err != nil {
			log.Println(err)
			continue
		}
		response.Message = response.Message + " " + respo
	}

	if marshaled, err := json.Marshal(response); err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(marshaled)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
func DoGetRequest(url string) (string, error) {
	log.Printf("processing get on %s\n", url)
	response, err := http.Get(url)
	responseBody := Response{}
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", err
	}
	return responseBody.Message, nil
}
