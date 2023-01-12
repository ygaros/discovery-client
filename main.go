package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"ygaros-discovery-client/client"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var serverClient client.Client

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Group(func(r chi.Router) {
		r.Get("/*", GetMapping)
		r.Get("/", GetMapping)
	})
	var err error
	serverClient, err = client.NewClientAndHeartBeat(
		"localhost",
		7654,
		"service-on-8000",
		"localhost",
		8000,
		false,
	)
	if err != nil {
		log.Println(err)
	}
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", 8000), r))

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
