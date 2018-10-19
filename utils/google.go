package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"ws-bigiot-ecoroutes/models"
)

var GOOGLE_DIRECTIONS_MODE_WALK = "walking"
var GOOGLE_DIRECTIONS_MODE_CYCLE = "bicycling"

type GoogleDirections struct {
	APIKey string `json:"key"`
}

var GOOGLE_DIRECTIONS_URL = "https://maps.googleapis.com/maps/api/directions/json?"

var GOOGLE_DIRECTIONS_STATUS_OK = "OK"

func (g *GoogleDirections) GetRoutingDistance(latini, lngini, latfin, lngfin float64, mode string) (*models.DirectionsReponse, error) {
	url := GOOGLE_DIRECTIONS_URL

	req, err := http.NewRequest("POST", url, nil)
	fmt.Println("url request:", req)

	arrayInitLocation := []string{fmt.Sprintf("%.6f", latini), fmt.Sprintf("%.6f", lngini)}
	arrayEndLocation := []string{fmt.Sprintf("%.6f", latfin), fmt.Sprintf("%.6f", lngfin)}

	extraParams := map[string]string{
		"origin":      strings.Join(arrayInitLocation, ","),
		"destination": strings.Join(arrayEndLocation, ","),
		"mode":        mode,
		"key":         g.APIKey,
	}
	fmt.Println("url params:", extraParams)
	values := addURLParameters(req.URL.Query(), extraParams)
	req.URL.RawQuery = values.Encode()
	fmt.Println("url request:", req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var response models.GoogleDirectionsResponse
	erro := json.Unmarshal(body, &response)
	if erro != nil {
		return nil, erro
	}

	var result models.DirectionsReponse
	result.Status = response.Status
	if response.Status == GOOGLE_DIRECTIONS_STATUS_OK {
		result.Distance = response.Routes[0].Legs[0].Distance.Value
		result.Duration = response.Routes[0].Legs[0].Duration.Value
	}

	return &result, nil
}

func addURLParameters(values url.Values, params map[string]string) url.Values {
	for key, value := range params {
		fmt.Println("url param key: ", key, " value: ", value)

		values.Add(key, value)
	}
	return values
}
