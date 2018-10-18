package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"bitbucket.org/sparsitytechnologies/go-ws-tools/models"
	"github.com/astaxie/beego"
)

func SetPathLoginKeys() {

	paths := [2]string{"PrivateKeyPath", "PublicKeyPath"}

	for _, element := range paths {
		path := beego.AppConfig.String(element)
		if _, err := os.Stat("vendor"); err == nil {
			path = "vendor/" + path
		} else {
			gopath := os.Getenv("GOPATH")
			path = gopath + "/src/" + path
		}
		beego.AppConfig.Set(strings.ToLower(element), path)

	}

}

func ProcessRequest(wResponseWriter http.ResponseWriter, rRequest *http.Request, requestObject interface{}) error {

	return ProcessRequestWithOptions(wResponseWriter, rRequest, requestObject, models.RequestOptions{})
}

func ProcessRequestWithOptions(wResponseWriter http.ResponseWriter, rRequest *http.Request, requestObject interface{}, requestOptions models.RequestOptions) error {

	var maxBytes = int64(1048576)
	if requestOptions.MaxBytes != 0 {
		maxBytes = requestOptions.MaxBytes
	}

	body, err := ioutil.ReadAll(io.LimitReader(rRequest.Body, maxBytes))
	if err != nil {
		//panic(err)
		//WriteResponse(wResponseWriter, nil, modelsTools.ResponseMessageTypeError, "["+modelsTools.ResponseMessageTypeError+"] "+err.Error())
		return err
	}
	if err := rRequest.Body.Close(); err != nil {
		//panic(err)
		//WriteResponse(wResponseWriter, nil, modelsTools.ResponseMessageTypeError, "["+modelsTools.ResponseMessageTypeError+"] "+err.Error())
		return err
	}

	wResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")

	erro := json.Unmarshal(body, &requestObject)

	return erro
}

func WriteResponse(wResponseWriter http.ResponseWriter, data interface{}, messageType, messageText string) {

	message := models.Message{
		Type: messageType,
		Text: messageText,
	}

	response := models.Response{
		Data:    data,
		Message: message,
	}

	wResponseWriter.WriteHeader(http.StatusOK)

	b, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(string(b))

	if _, err := wResponseWriter.Write(b); err != nil {
		panic(err)
	}
}

func IsAlive(wResponseWriter http.ResponseWriter, rRequest *http.Request, next http.HandlerFunc) {

	body, err := ioutil.ReadAll(io.LimitReader(rRequest.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := rRequest.Body.Close(); err != nil {
		panic(err)
	}
	fmt.Println(body)

	wResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")

	wResponseWriter.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(wResponseWriter).Encode("OK"); err != nil {
		panic(err)
		errorReturn(wResponseWriter, 500, err)
	}

}

func errorReturn(wResponseWriter http.ResponseWriter, code int, erro error) {
	wResponseWriter.WriteHeader(code) // unprocessable entity
	err := json.NewEncoder(wResponseWriter).Encode(erro)
	if err != nil {
		panic(err)
	}
}
