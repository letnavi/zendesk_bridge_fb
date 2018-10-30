package main

import (
	"encoding/json"
	"log"

	"github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/gin-gonic/gin"
)

//------------------------------------------------------------------------------------------------------------------
// ZENDESK
//------------------------------------------------------------------------------------------------------------------
// Create ticket, see doc
// https://developer.zendesk.com/rest_api/docs/core/tickets#create-ticket
//
// Create ticket with a set of parameters.
// Create comments if there is one in json
// To understand the features, see zendesk.Ticket obj
func createTicket(r []byte, c *gin.Context, z zendesk.Client) {

	ticketData := &zendesk.Ticket{}

	json.Unmarshal(r, &ticketData)

	z.CreateTicket(ticketData)

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

	_, err = z.UpdateTicket(*id, t)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}
