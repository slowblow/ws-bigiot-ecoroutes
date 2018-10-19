##Description:

Microservice to provide the best route between 2 points of BCN, providing a bikestation near to initial point and other bikestation near to end point.

##Dependencies

- package: github.com/astaxie/beego
  version: v1.10.1
- package: github.com/dgrijalva/jwt-go
  version: v3.2.0
- package: github.com/gorilla/mux
  version: v1.6.2
- package: github.com/kellydunn/golang-geo
  version: v0.7.0
- package: github.com/urfave/negroni
  version: v1.0.0

##Content:

###POST GetRoute:

#####Input parameters (JSON):

	* ini_point:
        * Latitude (float64)
		* Longitude (float64)
    * end_point
		* Latitude (float64)
		* Longitude (float64)

#####Output parameters:

	* Data: 
        JSON with data route

  	* Message
    	* Type: "SUCCESS", "ERROR"
    	* Text: description of error

#####Sample call:

```

curl -H "Content-Type: application/json" -d '{"ini_point": {"Latitude":41.38951,"Longitude":2.11295},"end_point": {"Latitude":41.35396877522612,"Longitude":2.1286922693252563}}' http://localhost:8438/ws-bigiot-ecoroutes/routes

```


#####Sample reponse:

```

{
    "Data": {
        "IniPoint": {
            "latitude": 41.387015399999996,
            "longitude": 2.1700471
        },
        "EndPoint": {
            "latitude": 41.375428299999996,
            "longitude": 2.1489419000000005
        },
        "IniBikeStation": {
            "slots": "9",
            "bikes": "12",
            "latitude": 41.387678,
            "longitude": 2.169587,
            "type": "BIKE",
            "status": "OPN",
            "distance_to_point": 0.08307661389495238
        },
        "EndBikeStation": {
            "slots": "26",
            "bikes": "5",
            "latitude": 41.376428,
            "longitude": 2.147734,
            "type": "BIKE",
            "status": "OPN",
            "distance_to_point": 0.15004936139260042
        },
        "DirectionsReponses": [
            {
                "Status": "OK",
                "Distance": 106,
                "Duration": 80
            },
            {
                "Status": "OK",
                "Distance": 2585,
                "Duration": 506
            },
            {
                "Status": "OK",
                "Distance": 106,
                "Duration": 80
            }
        ]
    },
    "Message": {
        "Type": "SUCCESS",
        "Text": "[SUCCESS] Response from BigIOT Services sucesfuly."
    }
}


```


##Docker

#####Build image

docker build --no-cache -f Dockerfile -t ws-bigiot-ecoroutes .


#####Run container

docker run -p 8438:8080 --name ws-bigiot-ecoroutes -d ws-bigiot-ecoroutes



