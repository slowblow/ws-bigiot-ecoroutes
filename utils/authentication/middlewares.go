package authentication

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"bitbucket.org/sparsitytechnologies/go-ws-tools/models"
	modelsTools "bitbucket.org/sparsitytechnologies/go-ws-tools/models"
	utilsTools "bitbucket.org/sparsitytechnologies/go-ws-tools/utils"
	httpWsTools "bitbucket.org/sparsitytechnologies/go-ws-tools/utils/http"
	modelUser "bitbucket.org/sparsitytechnologies/user-api/models"

	"github.com/astaxie/beego"
	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
)

func getCheckTokenWSLoginURL() string {
	var url = beego.AppConfig.String("CheckTokenWsLoginURL")
	fmt.Println("CheckTokenWsLoginURL:", url)
	return url
}

func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if getCheckTokenWSLoginURL() != "" {
		requireTokenAuthenticationHTTP(rw, req, next)
	} else {
		RequireTokenAuthenticationAPI(rw, req, next)
	}
}

/*
	Hacemos público el método para que el ws-login pueda usarlo directamente y así no entrar en un bucle infinito cuando está desplegado con la configuración compartida del deploy-conf ya que getCheckTokenWSLoginURL() siempre está informado en todos los modulos ya que comparten configuraciones.
*/
func RequireTokenAuthenticationAPI(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	token, err := getToken(rw, req)

	if err == nil && token.Valid {
		//fmt.Println("Token: ", token.Raw)
		if tokenTimeToLifeValidation(rw, token.Raw /*req,*/, next) {
			next(rw, req)
		}
	} else {
		fmt.Println("TOKEN: it isn't valid.", err.Error())
		models.PrintErrorResponse(rw, http.StatusUnauthorized, "TOKEN: it isn't valid. "+err.Error())
	}
}

func requireTokenAuthenticationHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	var authorization = req.Header.Get("Authorization")

	if authorization == "" {
		fmt.Println("TOKEN: No authorization in HEADER.")
		models.PrintErrorResponse(rw, http.StatusUnauthorized, "TOKEN: No authorization in HEADER.")
		return
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		fmt.Println("TOKEN: No Bearer in HEADER.")
		models.PrintErrorResponse(rw, http.StatusUnauthorized, "TOKEN: No Bearer in HEADER.")
		return
	}

	if len(authorization) < 8 {
		fmt.Println("TOKEN: No Token in HEADER.")
		models.PrintErrorResponse(rw, http.StatusUnauthorized, "TOKEN: No Token in HEADER.")
		return
	}

	token := authorization[7:len(authorization)]

	response, err := httpWsTools.HttpClientDoRequestWithToken("POST", getCheckTokenWSLoginURL(), token)

	if err != nil {
		fmt.Println(err.Error())
		models.PrintErrorResponse(rw, http.StatusUnauthorized, err.Error())
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(response.Body, 1048576))

	if err != nil {
		fmt.Println(err.Error())
		models.PrintErrorResponse(rw, http.StatusUnauthorized, err.Error())
		return
	}

	if err := response.Body.Close(); err != nil {
		fmt.Println(err.Error())
		models.PrintErrorResponse(rw, http.StatusUnauthorized, err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.Unmarshal(body, &models.Response{}); err != nil {
		var message string
		var e models.Error
		if erro := json.Unmarshal(body, &e); erro != nil {
			message = erro.Error()
		} else {
			message = e.Message
		}

		models.PrintErrorResponse(rw, http.StatusUnprocessableEntity, message)

		return
	}

	next(rw, req)

}

func getToken(rw http.ResponseWriter, req *http.Request) (*jwt.Token, error) {
	authBackend := InitJWTAuthenticationBackend()

	token, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return authBackend.PublicKey, nil
	})

	return token, err
}

func tokenTimeToLifeValidation(rw http.ResponseWriter, token string /*req *http.Request,*/, next http.HandlerFunc) (boolean bool) {

	//var authorization = req.Header.Get("Authorization")
	//var token = authorization[7:len(authorization)]
	//fmt.Println("Token " + token)

	userLoginDB := modelUser.FindUserLoginByToken(token)
	if userLoginDB.ID != 0 {
		var updatedAt = userLoginDB.UpdatedAt
		var createdAt = userLoginDB.CreatedAt

		var tokenTimeToLife, err1 = beego.AppConfig.Int("TokenTimeToLife")
		if err1 != nil {
			tokenTimeToLife = 30
		}
		updatedAt = updatedAt.Add(time.Duration(tokenTimeToLife) * time.Minute)

		var tokenMaxTimeToLife, err2 = beego.AppConfig.Int("TokenMaxTimeToLife")
		if err2 != nil {
			tokenMaxTimeToLife = 240
		}
		createdAt = updatedAt.Add(time.Duration(tokenMaxTimeToLife) * time.Minute)

		now := time.Now()

		if updatedAt.Before(now) {
			fmt.Println("TOKEN: Time exceded.")
			models.PrintErrorResponse(rw, http.StatusUnauthorized, "TOKEN: Time exceded.")
			return false
		} else if createdAt.Before(now) {
			fmt.Println("TOKEN: Max Time exceded.")
			models.PrintErrorResponse(rw, http.StatusUnauthorized, "TOKEN: Max Time exceded.")
			return false
		} else {
			modelUser.UpdateUserLogin(userLoginDB)
		}

	} else {
		fmt.Println("TOKEN: it doesn't exist.")
		models.PrintErrorResponse(rw, http.StatusUnauthorized, "TOKEN: it doesn't exist.")
		return false
	}

	return true
}

func RequireBasicAuthenticationHttp(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		//http.Error(rw, "authorization failed", http.StatusUnauthorized)

		fmt.Println("Security isn't a Basic Authorization")
		utilsTools.WriteResponse(rw, nil, modelsTools.ResponseMessageTypeError, "["+modelsTools.ResponseMessageTypeError+"] Authorization failed")

		return
	}

	payload, error := base64.StdEncoding.DecodeString(auth[1])
	if error != nil {
		fmt.Println("Security isn't a Basic Authorization (Decoding)")
		utilsTools.WriteResponse(rw, nil, modelsTools.ResponseMessageTypeError, "["+modelsTools.ResponseMessageTypeError+"] Authorization failed")
		return
	}

	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 || !validateBasicAuthenticationHttp(pair[0], pair[1]) {
		//http.Error(rw, "authorization failed", http.StatusUnauthorized)

		fmt.Println("Security isn't a valid Basic Authorization")
		utilsTools.WriteResponse(rw, nil, modelsTools.ResponseMessageTypeError, "["+modelsTools.ResponseMessageTypeError+"] Authorization failed")
		return
	}

	next(rw, req)
}

func validateBasicAuthenticationHttp(username, password string) bool {
	configUsername := GetServerBasicAuthenticationHttpUsername()

	configPassword := GetServerBasicAuthenticationHttpPassword()

	if configUsername != "" && configUsername == username && configPassword != "" && configPassword == password {
		return true
	}

	return false
}

func GetServerBasicAuthenticationHttpUsername() string {
	username := os.Getenv("SERVER_BASIC_AUTHENTICATION_HTTP_USERNAME")
	if username == "" {
		username = beego.AppConfig.String("ServerBasicAuthenticationHttpUsername")
	}

	return username
}

func GetServerBasicAuthenticationHttpPassword() string {
	password := os.Getenv("SERVER_BASIC_AUTHENTICATION_HTTP_PASSWORD")
	if password == "" {
		password = beego.AppConfig.String("ServerBasicAuthenticationHttpPassword")
	}

	return password
}
