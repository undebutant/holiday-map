package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Consts definitions
const JSON_PATH = "./data/data.json"

// Structs definitions
type Data struct {
	Markers    []Marker `json:"markers"`
	PhotoCount int      `json:"photoCount"`
}

type Marker struct {
	Name      string  `json:"name"`
	Latitude  string  `json:"latitude"`
	Longitude string  `json:"longitude"`
	Photos    []Photo `json:"photos"`
}

type Photo struct {
	Id          int    `json:"id"`
	FileName    string `json:"fileName"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

func getJsonData(context *gin.Context) Data {
	jsonData, fileError := os.Open(JSON_PATH)
	var data Data

	if fileError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
	}
	defer jsonData.Close()

	byteValue, decodingError := ioutil.ReadAll(jsonData)
	if decodingError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
	}
	json.Unmarshal(byteValue, &data)
	return data
}

func setJsonData(context *gin.Context, data Data) {
	file, marshalError := json.MarshalIndent(data, "", "    ")
	if marshalError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
	}
	writeError := ioutil.WriteFile(JSON_PATH, file, 0777)
	if writeError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
	}
}

// Route functions
func getMarkers(context *gin.Context) {
	data := getJsonData(context)
	context.JSON(http.StatusOK, data.Markers)
}

//WARNING : latitude and longitude strings must be trimmed in client (no unnecessary 0 at the end)
func getMarker(context *gin.Context) {
	data := getJsonData(context)
	latitude := context.Query("latitude")
	longitude := context.Query("longitude")
	for _, marker := range data.Markers {
		if marker.Latitude == latitude && marker.Longitude == longitude {
			context.JSON(http.StatusOK, marker)
			return
		}
	}
}

func addMarker(context *gin.Context) {
	var markerData Marker
	context.BindJSON(&markerData)
	currentData := getJsonData(context)
	currentData.Markers = append(currentData.Markers, markerData)
	setJsonData(context, currentData)
}

func main() {
	// Create a gin router with default middleware
	router := gin.Default()

	// Serve static web files
	router.StaticFile("/", "./resources/main.html")
	router.StaticFile("/static/main.css", "./resources/main.css")

	// Serve photos
	router.Static("/photos", "./data/photos")

	// Dynamic routing
	router.GET("/markers", getMarkers)
	router.GET("/marker", getMarker)
	router.PUT("/marker", addMarker)

	// Listen and serve on localhost:8080
	router.Run()
}
