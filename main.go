package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//--------------------------------------------------------------------------------------------------------------------//
// Consts definitions
//--------------------------------------------------------------------------------------------------------------------//
const JSON_PATH = "./data/data.json"

//--------------------------------------------------------------------------------------------------------------------//
// Structs definitions
//--------------------------------------------------------------------------------------------------------------------//
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

//--------------------------------------------------------------------------------------------------------------------//
// Utilities functions
//--------------------------------------------------------------------------------------------------------------------//
func getJsonData(context *gin.Context) Data {
	var data Data

	jsonData, fileError := os.Open(JSON_PATH)
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

func removeMarkerFromArray(markerArray []Marker, removeIndex int) []Marker {
	markerArray[removeIndex] = markerArray[len(markerArray)-1]
	return markerArray[:len(markerArray)-1]
}

//--------------------------------------------------------------------------------------------------------------------//
// Route functions
//--------------------------------------------------------------------------------------------------------------------//
func getMarkers(context *gin.Context) {
	data := getJsonData(context)
	context.JSON(http.StatusOK, data.Markers)
}

// WARNING: latitude and longitude strings must be trimmed client side (no unnecessary 0 at the end)
func getMarker(context *gin.Context) {
	data := getJsonData(context)
	latitude := context.Param("latitude")
	longitude := context.Param("longitude")
	for _, marker := range data.Markers {
		if marker.Latitude == latitude && marker.Longitude == longitude {
			context.JSON(http.StatusOK, marker)
			return
		}
	}
	context.AbortWithStatus(http.StatusNotFound)
}

func addMarker(context *gin.Context) {
	var markerData Marker
	context.BindJSON(&markerData)
	currentData := getJsonData(context)
	currentData.Markers = append(currentData.Markers, markerData)
	setJsonData(context, currentData)
}

func editMarker(context *gin.Context) {
	data := getJsonData(context)
	latitude := context.Param("latitude")
	longitude := context.Param("longitude")
	for index, marker := range data.Markers {
		if marker.Latitude == latitude && marker.Longitude == longitude {
			var markerData Marker
			context.BindJSON(&markerData)

			marker.Latitude = markerData.Latitude
			marker.Longitude = markerData.Longitude
			marker.Name = markerData.Name

			data.Markers[index] = marker

			setJsonData(context, data)
			context.Status(http.StatusOK)
			return
		}
	}
	context.AbortWithStatus(http.StatusNotFound)
}

func deleteMarker(context *gin.Context) {
	data := getJsonData(context)
	latitude := context.Param("latitude")
	longitude := context.Param("longitude")
	for index, marker := range data.Markers {
		if marker.Latitude == latitude && marker.Longitude == longitude {
			data.Markers = removeMarkerFromArray(data.Markers, index)
			setJsonData(context, data)
			context.Status(http.StatusOK)
			return
		}
	}
	context.AbortWithStatus(http.StatusNotFound)
}

//--------------------------------------------------------------------------------------------------------------------//
// Main
//--------------------------------------------------------------------------------------------------------------------//
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
	router.GET("/marker/:latitude/:longitude", getMarker)
	router.PUT("/marker", addMarker)
	router.POST("/marker/:latitude/:longitude", editMarker)
	router.DELETE("/marker/:latitude/:longitude", deleteMarker)

	// Listen and serve on localhost:8080
	router.Run()
}
