package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

//--------------------------------------------------------------------------------------------------------------------//
// Consts definitions
//--------------------------------------------------------------------------------------------------------------------//
const JSON_PATH = "./data/data.json"
const PHOTOS_PATH = "./data/photos/"

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
	Date        string `json:"date"` // Format YYYY-MM-DD
}

type MinimalMarkers struct {
	Markers []MinimalMarker `json:"markers"`
}

type MinimalMarker struct {
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

//--------------------------------------------------------------------------------------------------------------------//
// Utilities functions
//--------------------------------------------------------------------------------------------------------------------//
func getJsonData(context *gin.Context) (Data, error) {
	var data Data

	jsonData, fileError := os.Open(JSON_PATH)
	if fileError != nil {
		return data, fileError
	}
	defer jsonData.Close()

	byteValue, decodingError := ioutil.ReadAll(jsonData)
	if decodingError != nil {
		return data, decodingError
	}

	json.Unmarshal(byteValue, &data)
	return data, nil
}

func setJsonData(context *gin.Context, data Data) error {
	file, marshalError := json.MarshalIndent(data, "", "    ")
	if marshalError != nil {
		return marshalError
	}

	writeError := ioutil.WriteFile(JSON_PATH, file, 0777)
	if writeError != nil {
		return writeError
	}

	return nil
}

func removeMarkerFromArray(markerArray []Marker, removeIndex int) []Marker {
	if len(markerArray) == 1 {
		markerArray = markerArray[:0]
	} else {
		markerArray[removeIndex] = markerArray[len(markerArray)-1]
		markerArray = markerArray[:len(markerArray)-1]
	}

	return markerArray
}

func removePhotoFromArray(photoArray []Photo, removeIndex int) []Photo {
	if len(photoArray) == 1 {
		photoArray = photoArray[:0]
	} else {
		photoArray[removeIndex] = photoArray[len(photoArray)-1]
		photoArray = photoArray[:len(photoArray)-1]
	}

	return photoArray
}

//--------------------------------------------------------------------------------------------------------------------//
// Route functions
//--------------------------------------------------------------------------------------------------------------------//
func getMarkers(context *gin.Context) {
	data, readError := getJsonData(context)
	if readError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Remove photos and photoCount from response
	var dataToSend MinimalMarkers
	for _, marker := range data.Markers {
		var minimalMarker MinimalMarker
		minimalMarker.Name = marker.Name
		minimalMarker.Latitude = marker.Latitude
		minimalMarker.Longitude = marker.Longitude

		dataToSend.Markers = append(dataToSend.Markers, minimalMarker)
	}

	context.JSON(http.StatusOK, data.Markers)
}

// WARNING - latitude and longitude strings must be trimmed on client side (no unnecessary 0 at the end)
func getMarker(context *gin.Context) {
	data, readError := getJsonData(context)
	if readError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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

	currentData, readError := getJsonData(context)
	if readError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	currentData.Markers = append(currentData.Markers, markerData)
	writeError := setJsonData(context, currentData)
	if writeError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusOK)
}

func editMarker(context *gin.Context) {
	data, readError := getJsonData(context)
	if readError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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
			writeError := setJsonData(context, data)
			if writeError != nil {
				context.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			context.Status(http.StatusOK)
			return
		}
	}

	context.AbortWithStatus(http.StatusNotFound)
}

func deleteMarker(context *gin.Context) {
	data, readError := getJsonData(context)
	if readError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	latitude := context.Param("latitude")
	longitude := context.Param("longitude")
	for index, marker := range data.Markers {
		if marker.Latitude == latitude && marker.Longitude == longitude {
			// Store path to related photos
			var photosToDelete []string
			for _, photo := range marker.Photos {
				photosToDelete = append(photosToDelete, PHOTOS_PATH+photo.FileName)
			}

			// Remove marker from data
			data.Markers = removeMarkerFromArray(data.Markers, index)
			writeError := setJsonData(context, data)
			if writeError != nil {
				context.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			// Remove related photos
			for _, photoPath := range photosToDelete {
				deletionError := os.Remove(photoPath)
				if deletionError != nil {
					fmt.Println("Error deleting " + photoPath + " with '" + deletionError.Error() + "'")
				}
			}

			context.Status(http.StatusOK)
			return
		}
	}

	context.AbortWithStatus(http.StatusNotFound)
}

func addPhoto(context *gin.Context) {
	data, readError := getJsonData(context)
	if readError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	file, fileError := context.FormFile("photo")
	if fileError != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	latitude := context.Param("latitude")
	longitude := context.Param("longitude")
	for index, marker := range data.Markers {
		if marker.Latitude == latitude && marker.Longitude == longitude {
			data.PhotoCount++

			var photo Photo
			photo.Id = data.PhotoCount
			photo.FileName = strconv.Itoa(data.PhotoCount) + "_" + context.Query("fileName")
			photo.Description = context.Query("description")
			photo.Date = context.Query("date")

			marker.Photos = append(marker.Photos, photo)

			// Upload the file to the specified destination
			context.SaveUploadedFile(file, PHOTOS_PATH+photo.FileName)

			data.Markers[index] = marker
			writeError := setJsonData(context, data)
			if writeError != nil {
				context.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			context.Status(http.StatusOK)
			return
		}
	}

	context.AbortWithStatus(http.StatusNotFound)
}

func editPhoto(context *gin.Context) {
	data, readError := getJsonData(context)
	if readError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	latitude := context.Param("latitude")
	longitude := context.Param("longitude")
	photoId, _ := strconv.Atoi(context.Param("photoId"))
	for index, marker := range data.Markers {
		if marker.Latitude == latitude && marker.Longitude == longitude {
			for photoIndex, photo := range marker.Photos {
				if photo.Id == photoId {
					var photoData Photo
					context.BindJSON(&photoData)

					photo.Description = photoData.Description
					photo.Date = photoData.Date

					data.Markers[index].Photos[photoIndex] = photo
					writeError := setJsonData(context, data)
					if writeError != nil {
						context.AbortWithStatus(http.StatusInternalServerError)
						return
					}

					context.Status(http.StatusOK)
					return
				}
			}
		}
	}

	context.AbortWithStatus(http.StatusNotFound)
}

func deletePhoto(context *gin.Context) {
	data, readError := getJsonData(context)
	if readError != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	latitude := context.Param("latitude")
	longitude := context.Param("longitude")
	photoId, _ := strconv.Atoi(context.Param("photoId"))
	for index, marker := range data.Markers {
		if marker.Latitude == latitude && marker.Longitude == longitude {
			for photoIndex, photo := range marker.Photos {
				if photo.Id == photoId {
					pathToDelete := PHOTOS_PATH + photo.FileName

					// Remove photo from data
					marker.Photos = removePhotoFromArray(marker.Photos, photoIndex)

					data.Markers[index] = marker
					writeError := setJsonData(context, data)
					if writeError != nil {
						context.AbortWithStatus(http.StatusInternalServerError)
						return
					}

					// Remove photo file
					deletionError := os.Remove(pathToDelete)
					if deletionError != nil {
						fmt.Println("Error deleting " + pathToDelete + " with '" + deletionError.Error() + "'")
					}

					context.Status(http.StatusOK)
					return
				}
			}
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
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.StaticFile("/static/main.css", "./resources/main.css")

	// Serve photos
	router.Static("/photos", "./data/photos")

	// Dynamic routing
	router.GET("/markers", getMarkers)
	router.GET("/marker/:latitude/:longitude", getMarker)
	router.PUT("/marker", addMarker)
	router.POST("/marker/:latitude/:longitude", editMarker)
	router.DELETE("/marker/:latitude/:longitude", deleteMarker)

	router.POST("/marker/:latitude/:longitude/photo", addPhoto)
	router.POST("/marker/:latitude/:longitude/photo/:photoId", editPhoto)
	router.DELETE("/marker/:latitude/:longitude/photo/:photoId", deletePhoto)

	// Listen and serve on localhost:8080
	router.Run()
}
