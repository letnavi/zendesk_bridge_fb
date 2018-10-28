package main

import (
	"log"
	"net/http"

	"github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/gin-gonic/gin"
)

type Configuration struct {
	ZendeskDomain   []string
	ZendeskUsername []string
	ZendeskPassword []string
}

const jsonData = "resources/create.json"

func main() {

	r := gin.Default()
	r.Use(LiberalCORS)
	r.GET("/ping", pong)
	r.POST("/create", create)
	r.Run(":9090")
}

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func LiberalCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	if c.Request.Method == "OPTIONS" {
		if len(c.Request.Header["Access-Control-Request-Headers"]) > 0 {
			c.Header("Access-Control-Allow-Headers", c.Request.Header["Access-Control-Request-Headers"][0])
		}
		c.AbortWithStatus(http.StatusOK)
	}
}

func create(c *gin.Context) {

	// Коннескт
	client, err := zendesk.NewClient("", "", "")

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"status": action.createAction(jsonData, client),
	})
}
