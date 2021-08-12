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
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Photos    []Photo `json:"photos"`
}

type Photo struct {
	Id          int    `json:"id"`
	FileName    string `json:"fileName"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

// Route functions
func getMarkers(context *gin.Context) {
	jsonData, fileError := os.Open(JSON_PATH)

	if fileError != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	defer jsonData.Close()

	byteValue, decodingError := ioutil.ReadAll(jsonData)
	if decodingError != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	var data Data
	json.Unmarshal(byteValue, &data)

	context.JSON(http.StatusOK, data.Markers)
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

	// Listen and serve on localhost:8080
	router.Run()
}
