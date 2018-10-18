package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func PrintErrorResponse(wResponseWriter http.ResponseWriter, errorCode int, errorMessage string) {

	if wResponseWriter != nil {
		var er Error
		er.Code = errorCode
		er.Message = errorMessage

		wResponseWriter.WriteHeader(er.Code)

		b, err := json.Marshal(er)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))

		if _, err := wResponseWriter.Write(b); err != nil {
			panic(err)
		}
	} else {
		panic(strconv.Itoa(errorCode) + " - " + errorMessage)
	}

}
