package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
	"ygaros-discovery-client/client"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Group(func(r chi.Router) {
		r.Get("/", GetMapping)
		r.Get("/*", GetMapping)
	})
	err := client.NewClientAndHeartBeat(
		"localhost",
		7654,
		"service-on-8000",
		"localhost",
		8000,
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
		Message: "hello from service on 8000",
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
