package models

type DurationResponse struct {
	Text  string `json:"text"`
	Value int64  `json:"value"`
}

type DistanceResponse struct {
	Text  string `json:"text"`
	Value int64  `json:"value"`
}

type LegResponse struct {
	Distance DistanceResponse `json:"distance"`
	Duration DurationResponse `json:"duration"`
}

type RoutesResponse struct {
	Legs []LegResponse `json:"legs"`
}

type GoogleDirectionsResponse struct {
	Status string           `json:"status"`
	Routes []RoutesResponse `json:"routes"`
}

type DirectionsReponse struct {
	Status   string `valid:"-"`
	Distance int64  `valid:"-"`
	Duration int64  `valid:"-"`
}
