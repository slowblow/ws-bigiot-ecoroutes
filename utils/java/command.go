package java

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"ws-bigiot-services/models"

	"github.com/astaxie/beego"
)

func Command() ([]models.BikeStation, error) {
	var bikeStations []models.BikeStation

	var file []byte
	var err error

	fake := beego.AppConfig.String("Fake")
	if fake == "" || fake == "false" {
		// delete file
		err := os.Remove("out.json")
		if err != nil {
			fmt.Println("==> error deleting file: ", err.Error())
		} else {
			fmt.Println("==> done deleting file")
		}

		path, _ := exec.LookPath("java")
		fmt.Println("========== START JAR ========== \n ")
		cmd := exec.Command(path, "-cp", "java-example-consumer.jar:.", "org.bigiot.examples.ExampleConsumer")

		_, err = cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("==========  ERROR CMD JAR ========== \n ")
			return bikeStations, err
		}
		//fmt.Println(string(out))
		fmt.Println("==========  END  JAR ========== \n ")

		file, err = ioutil.ReadFile("out.json")
		if err == nil {
			err = json.Unmarshal(file, &bikeStations)

			if err != nil {
				file, err = ioutil.ReadFile("harcoded.json")
				if err == nil {
					err = json.Unmarshal(file, &bikeStations)
				}
			}
		}
	} else {
		file, err = ioutil.ReadFile("harcoded.json")
		if err == nil {
			err = json.Unmarshal(file, &bikeStations)
		}
	}

	return bikeStations, err
}

func parseFileResponse(fileResponsePath string) (string, error) {

	var fileJSON string

	b, err := ioutil.ReadFile(fileResponsePath) // just pass the file name
	if err == nil {
		fileJSON = string(b)
	}

	return fileJSON, err
}
