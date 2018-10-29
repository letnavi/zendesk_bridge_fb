package main

import (
	"encoding/json"
	"github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// port for server
const port = ":9090"

// Processing facebook and zendesk requests
func main() {

	// zendesk client
	// sub-domain, email/login and password
	client, err := zendesk.NewClient("", "", "")

	// if not connect
	if err != nil {
		log.Fatal(err)
	}

	// gin route
	route := gin.Default()
	r := route.Group("/api/v1")
	r.Use(LiberalCORS)

	// test ping pong
	r.GET("/ping", pong)

	//------------------------------------------------------------------------------------------------------------------
	// FACEBOOK
	//------------------------------------------------------------------------------------------------------------------

	// add comment for facebook
	r.POST("/fb/comment", func(c *gin.Context) {
		createWorkplaceComment(c, client)
	})

	//------------------------------------------------------------------------------------------------------------------
	// ZENDESK
	//------------------------------------------------------------------------------------------------------------------

	// add ticket
	r.POST("/ticket", func(c *gin.Context) {
		createTicket(c, client)
	})

	// update ticket and/or add comment
	r.PUT("/ticket", func(c *gin.Context) {
		updateTicket(c, client)
	})
	r.Run(port)
}

func createWorkplaceComment(c *gin.Context, client zendesk.Client) {

}

// Test func
// Ping Pong
func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//
func LiberalCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	if c.Request.Method == "OPTIONS" {
		if len(c.Request.Header["Access-Control-Request-Headers"]) > 0 {
			c.Header("Access-Control-Allow-Headers", c.Request.Header["Access-Control-Request-Headers"][0])
		}
		c.AbortWithStatus(http.StatusOK)
	}
}

// Create ticket, see doc
// https://developer.zendesk.com/rest_api/docs/core/tickets#create-ticket
//
// Create ticket with a set of parameters.
// Create comments if there is one in json
// To understand the features, see zendesk.Ticket obj
func createTicket(c *gin.Context, z zendesk.Client) {

	ticketData := &zendesk.Ticket{}
	body, err := c.GetRawData()

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &ticketData)

	_, err = z.CreateTicket(ticketData)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

// Update ticket, see doc
// https://developer.zendesk.com/rest_api/docs/core/tickets#update-ticket
//
// Ticket update
// Add and update comment
// Depending on json content
func updateTicket(c *gin.Context, z zendesk.Client) {

	t := &zendesk.Ticket{}
	body, err := c.GetRawData()

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &t)

	id := t.ID

	_, err = z.UpdateTicket(id, t)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}
