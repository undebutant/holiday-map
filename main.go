package main

import (
    //"net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    // Creates a gin router with default middleware
    router := gin.Default()
    
    // Serves static files
    router.StaticFile("/", "./resources/main.html")
    router.StaticFile("/static/main.css", "./resources/main.css")

    // Listen and serves on localhost:8080
    router.Run()
}
