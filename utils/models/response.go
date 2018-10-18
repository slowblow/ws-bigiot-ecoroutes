package models

const (
	ResponseMessageTypeError   = "ERROR"
	ResponseMessageTypeWarning = "WARNING"
	ResponseMessageTypeSuccess = "SUCCESS"
)

type Response struct {
	Data    interface{}
	Message Message
}

type Message struct {
	Type string
	Text string
}

type WSResponse struct {
	Status       Status        `json:"status"`
	ErrorDetails *ErrorDetails `json:"error_details,omitemptys"`
	Content      interface{}   `json:"content,omitemptys"`
}

type ErrorDetails struct {
	Code      int    `json:"code"`
	Exception string `json:"exception"`
	Message   string `json:"message"`
}

type Status struct {
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
}

type BikeStation struct {
	Slots           int     `json:"slots,string"`
	Bikes           int     `json:"bikes,string"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Type            string  `json:"type"`
	Status          string  `json:"status"`
	DistanceToPoint float64 `json:"distance_to_point"`
}

type CloserBikeStations struct {
	IniPointList []BikeStation
	EndPointList []BikeStation
}

type Route struct {
	IniPoint           Point
	EndPoint           Point
	IniBikeStation     BikeStation
	EndBikeStation     BikeStation
	DirectionsReponses []DirectionsReponse
}
