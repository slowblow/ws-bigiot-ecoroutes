package services

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sort"

	"ws-bigiot-services/models"
	"ws-bigiot-services/utils"

	"ws-bigiot-services/utils/java"

	geo "github.com/kellydunn/golang-geo"
)

func GetRoutes(wResponseWriter http.ResponseWriter, rRequest *http.Request, next http.HandlerFunc) {

	var request models.Request

	body, err := ioutil.ReadAll(io.LimitReader(rRequest.Body, 1048576))

	if err != nil {
		utils.WriteResponse(wResponseWriter, nil, models.ResponseMessageTypeError, "["+models.ResponseMessageTypeError+"] "+err.Error())

		return
	}

	if err := rRequest.Body.Close(); err != nil {
		utils.WriteResponse(wResponseWriter, nil, models.ResponseMessageTypeError, "["+models.ResponseMessageTypeError+"] "+err.Error())
		return
	}

	wResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.Unmarshal(body, &request); err != nil {
		utils.WriteResponse(wResponseWriter, nil, models.ResponseMessageTypeError, "["+models.ResponseMessageTypeError+"] "+err.Error())
	} else {

		wResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")

		response, err := BigIOTServicesAPI(request.IniPoint, request.EndPoint)

		if err != nil {
			utils.WriteResponse(wResponseWriter, nil, models.ResponseMessageTypeError, "["+models.ResponseMessageTypeError+"] "+err.Error())
		} else {

			googleResponse, err := GoogleServicesAPI(request.IniPoint, request.EndPoint, response)

			sort.Sort(ByTime(googleResponse))
			if err != nil {
				utils.WriteResponse(wResponseWriter, nil, models.ResponseMessageTypeError, "["+models.ResponseMessageTypeError+"] "+err.Error())
			} else {
				utils.WriteResponse(wResponseWriter, googleResponse[0], models.ResponseMessageTypeSuccess, "["+models.ResponseMessageTypeSuccess+"] Response from BigIOT Services sucesfuly.")
			}
		}
	}
}

func BigIOTServicesAPI(initPoint models.Point, endPoint models.Point) (models.CloserBikeStations, error) {
	closerBikeStations := models.CloserBikeStations{}

	bikeStations, err := java.Command()

	if err != nil {

		return closerBikeStations, err
	}

	closerBikeStations = findCloserBikeStations(initPoint, endPoint, bikeStations)

	return closerBikeStations, err
}

func findCloserBikeStations(iniPoint models.Point, endPoint models.Point, bikeStations []models.BikeStation) models.CloserBikeStations {
	reducedIniPointList := []models.BikeStation{}
	reducedEndPointList := []models.BikeStation{}
	closerBikeStations := models.CloserBikeStations{}

	sortedBikeStationsByIniPoint := sortBikeStations(iniPoint, bikeStations)
	sortedBikeStationsByEndPoint := sortBikeStations(endPoint, bikeStations)

	for _, item := range sortedBikeStationsByIniPoint {
		if len(reducedIniPointList) < 3 && item.Bikes > 0 && item.Type == "BIKE" {
			reducedIniPointList = append(reducedIniPointList, item)
		}
	}
	closerBikeStations.IniPointList = reducedIniPointList

	for _, item := range sortedBikeStationsByEndPoint {
		if len(reducedEndPointList) < 3 && item.Bikes > 0 && item.Type == "BIKE" {
			reducedEndPointList = append(reducedEndPointList, item)
		}
	}
	closerBikeStations.EndPointList = reducedEndPointList

	pointsToShow := []models.Point{}

	pointsToShow = append(pointsToShow, iniPoint)

	for _, bikeStation := range reducedIniPointList {
		point := models.Point{
			Latitude:  bikeStation.Latitude,
			Longitude: bikeStation.Longitude,
		}
		pointsToShow = append(pointsToShow, point)

	}

	pointsToShow = append(pointsToShow, endPoint)

	for _, bikeStation := range reducedEndPointList {
		point := models.Point{
			Latitude:  bikeStation.Latitude,
			Longitude: bikeStation.Longitude,
		}
		pointsToShow = append(pointsToShow, point)

	}

	utils.SaveGeojson(pointsToShow)

	return closerBikeStations
}

func sortBikeStations(pointToCheck models.Point, bikeStations []models.BikeStation) []models.BikeStation {
	bikeStationsResult := []models.BikeStation{}
	pToCheck := geo.NewPoint(pointToCheck.Latitude, pointToCheck.Longitude)
	for _, bikeStation := range bikeStations {
		p := geo.NewPoint(bikeStation.Latitude, bikeStation.Longitude)

		bikeStation.DistanceToPoint = pToCheck.GreatCircleDistance(p)
		bikeStationsResult = append(bikeStationsResult, bikeStation)
	}

	//sort.Slice(bikeStationsResult, func(i, j int) bool { return bikeStations[i].DistanceToPoint < bikeStations[j].DistanceToPoint })
	sort.Sort(ByCloserToPoint(bikeStationsResult))

	return bikeStationsResult
}

type ByCloserToPoint []models.BikeStation

func (lista ByCloserToPoint) Len() int {
	return len(lista)
}

func (lista ByCloserToPoint) Swap(i, j int) {
	lista[i], lista[j] = lista[j], lista[i]
}

func (lista ByCloserToPoint) Less(i, j int) bool {
	return lista[i].DistanceToPoint < lista[j].DistanceToPoint
}

type ByTime []models.Route

func (lista ByTime) Len() int {
	return len(lista)
}

func (lista ByTime) Swap(i, j int) {
	lista[i], lista[j] = lista[j], lista[i]
}

func (lista ByTime) Less(i, j int) bool {
	time_i := lista[i].DirectionsReponses[0].Duration
	time_i += lista[i].DirectionsReponses[1].Duration
	time_i += lista[i].DirectionsReponses[2].Duration

	time_j := lista[j].DirectionsReponses[0].Duration
	time_j += lista[j].DirectionsReponses[1].Duration
	time_j += lista[j].DirectionsReponses[2].Duration
	return time_i < time_j
}

func GoogleServicesAPI(initPoint models.Point, endPoint models.Point, closerBikeStations models.CloserBikeStations) ([]models.Route, error) {
	routes := []models.Route{}
	var err error
	googleDirections := &utils.GoogleDirections{
		APIKey: "AIzaSyATqqadxqycHK73sIH7tl3xv7IfOiTndBA",
	}

	IniDirectionsReponses := []models.DirectionsReponse{}
	EndDirectionsReponses := []models.DirectionsReponse{}

	for _, ini := range closerBikeStations.IniPointList {
		directionsReponse := &models.DirectionsReponse{}
		directionsReponse, err = googleDirections.GetRoutingDistance(initPoint.Latitude, initPoint.Longitude, ini.Latitude, ini.Longitude, utils.GOOGLE_DIRECTIONS_MODE_WALK)
		if err != nil {
			return routes, err
		}
		IniDirectionsReponses = append(IniDirectionsReponses, *directionsReponse)
	}

	for _, end := range closerBikeStations.EndPointList {
		directionsReponse := &models.DirectionsReponse{}
		directionsReponse, err = googleDirections.GetRoutingDistance(endPoint.Latitude, endPoint.Longitude, end.Latitude, end.Longitude, utils.GOOGLE_DIRECTIONS_MODE_WALK)
		if err != nil {
			return routes, err
		}
		EndDirectionsReponses = append(EndDirectionsReponses, *directionsReponse)
	}

	for i, ini := range closerBikeStations.IniPointList {

		for j, end := range closerBikeStations.EndPointList {

			route := models.Route{}
			directionsReponse := &models.DirectionsReponse{}
			route.IniPoint = initPoint
			route.EndPoint = endPoint
			route.IniBikeStation = ini
			route.EndBikeStation = end
			directionsReponse, err = googleDirections.GetRoutingDistance(ini.Latitude, ini.Longitude, end.Latitude, end.Longitude, utils.GOOGLE_DIRECTIONS_MODE_CYCLE)

			route.DirectionsReponses = append(route.DirectionsReponses, IniDirectionsReponses[i])
			route.DirectionsReponses = append(route.DirectionsReponses, *directionsReponse)
			route.DirectionsReponses = append(route.DirectionsReponses, IniDirectionsReponses[j])

			if err != nil {
				return routes, err
			}

			routes = append(routes, route)
		}
	}

	return routes, nil
}
