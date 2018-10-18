package utils

import (
	"io"
	"os"
	"strconv"
	"strings"

	"ws-bigiot-ecoroutes/models"
)

const (
	geojson_head = "{" +
		"\"type\": \"FeatureCollection\"," +
		"\"features\": ["

	geojson_point_ini = "{" +
		"\"type\": \"Feature\"," +
		"\"properties\": {}," +
		"\"geometry\": {" +
		"\"type\": \"Point\"," +
		"\"coordinates\": ["

	geojson_point_fin = "]}}"
)

func SaveGeojson(points []models.Point) {
	/**/
	geojson := geojson_head

	//
	for i, point := range points {
		if i > 0 {
			geojson += ","
		}
		geojson += geojson_point_ini

		geojson += strconv.FormatFloat(point.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(point.Latitude, 'f', -1, 64)

		geojson += geojson_point_fin

	}
	geojson += "]}"

	WriteStringToFile("Points.json", geojson)
}

func WriteStringToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}

	return nil
}
