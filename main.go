package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/astaxie/beego"

	"ws-bigiot-ecoroutes/routers"
)

func main() {

	routers.ArrayRoutesWS = append(routers.ArrayRoutesWS, routers.RoutesWSArray...)

	router := routers.NewRouter()

	port := os.Getenv("PORT")

	if port == "" {
		port = beego.AppConfig.String("httpport")
	}

	fmt.Println("http server Running on internal PORT: ", port)

	log.Fatal(http.ListenAndServe(":"+port, router))

}
