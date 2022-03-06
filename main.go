package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	gkhttp "github.com/go-kit/kit/transport/http"
	"log"
	"net/http"
)

type TalkToMe interface {
	HowsTheWeather(string) (int, string, error)
}

type weatherRequest struct {
	ZipCode string `json:"zip_code"`
}

type weatherResponse struct {
	Temperature int    `json:"temperature"`
	Description string `json:"description"`
}

type talkToMe struct{}

func (*talkToMe) HowsTheWeather(zipCode string) (int, string, error) {
	if zipCode == "23059" {
		return 90, "it's hot!", nil
	} else {
		return -1, "i don't know", fmt.Errorf("ZipCodeException")
	}
}

//Endpoints
//simple adapters to convert each of our service’s methods into an endpoint
//type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
func weatherEndpoint(svc TalkToMe) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(weatherRequest)
		temp, description, err := svc.HowsTheWeather(req.ZipCode)

		if err != nil {
			return weatherResponse{temp, description}, err
		}

		//response in json format
		return weatherResponse{temp, description}, nil
	}
}

//type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
func decodeWeatherRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request weatherRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main() {

	t := talkToMe{}

	//Transports to expose your service to the outside world
	weatherHandler := gkhttp.NewServer(
		weatherEndpoint(&t),
		decodeWeatherRequest,
		encodeResponse,
	)

	http.Handle("/weather", weatherHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
