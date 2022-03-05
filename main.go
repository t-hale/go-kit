package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	gkhttp "github.com/go-kit/kit/transport/http"
	"log"
	"net/http"
)

//type TalkToMe interface {
//	HowsTheWeather(string) string
//}

type weatherRequest struct {
	ZipCode string `json:"zip_code"`
}

type weatherResponse struct {
	Temperature int    `json:"temperature"`
	Description string `json:"description"`
	Err         string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type talkToMe struct{}

func (*talkToMe) HowsTheWeather(zipcode string) (int, error) {
	return 90, nil
}

//Endpoints
//simple adapters to convert each of our service’s methods into an endpoint
//type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
func weatherEndpoint(svc talkToMe) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(weatherRequest)
		v, err := svc.HowsTheWeather(req.ZipCode)
		if err != nil {
			//response in json format
			return weatherResponse{-1, "here's a description", "we had an error"}, nil
		}
		//response in json format
		return weatherResponse{v, "a cool description", ""}, nil
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
		weatherEndpoint(t),
		decodeWeatherRequest,
		encodeResponse,
	)

	http.Handle("/weather", weatherHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
