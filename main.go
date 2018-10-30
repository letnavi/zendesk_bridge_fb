package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Domain   string
	Username string
	Password string
	Port     string
	Token    string
}

type Feed struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Period string `json:"period"`
	Values []struct {
		Value int `json:"value"`
	} `json:"values"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

const (
	ConfigFile = "conf.toml"
	StatusOpen = "open"
)

var route = gin.Default()
var r = *route.Group("/api/v1")
var conf Config

func init() {
	gin.SetMode(gin.ReleaseMode)

	r.Use(LiberalCORS)
	// init config
	if _, err := toml.DecodeFile(ConfigFile, &conf); err != nil {
		log.Fatal(err)
	}
}

// Processing facebook and zendesk requests
func main() {

	// zendesk client
	// sub-domain, email/login and password
	client, err := zendesk.NewClient(conf.Domain, conf.Username, conf.Password)
	// if not connect
	if err != nil {
		log.Fatal(err)
	}

	// test ping pong
	r.GET("/ping", pong)
	// verify token fb
	r.GET("/fb", func(c *gin.Context) {
		Verify(conf.Token, c.Writer, c.Request)
	})
	// add comment for facebook
	r.POST("/fb", func(c *gin.Context) {
		createWorkplaceComment(c, client)
	})
	// add ticket
	r.POST("/ticket", func(c *gin.Context) {
		if req, err := toZendesk(c); err != nil {
			log.Fatal("что то пошло не так")
		} else {
			createTicket(req, c, client)
		}
	})
	// update ticket and/or add comment
	r.PUT("/ticket", func(c *gin.Context) {
		updateTicket(c, client)
	})

	route.Run(conf.Port)
}

// convert request fb to zendesk
// create zendesk ticket
func toZendesk(c *gin.Context) (zendesk.Ticket, error) {
	body, err := c.GetRawData()
	if err != nil {
		log.Fatal(err)
	}

	feed := Feed{}
	if err = json.Unmarshal(body, &feed); err != nil {
		log.Fatal(err)
	}

	status := StatusOpen
	ticket := zendesk.Ticket{
		Subject:     &feed.Title,
		Description: &feed.Description,
		Status:      &status,
	}

	return ticket, nil
}

//------------------------------------------------------------------------------------------------------------------
// ZENDESK
//------------------------------------------------------------------------------------------------------------------

// Create ticket, see doc
// https://developer.zendesk.com/rest_api/docs/core/tickets#create-ticket
//
// Create ticket with a set of parameters.
// Create comments if there is one in json
// To understand the features, see zendesk.Ticket obj
func createTicket(t zendesk.Ticket, c *gin.Context, z zendesk.Client) {
	if ticket, err := z.CreateTicket(&t); err != nil {
		c.JSON(200, gin.H{
			"status": "ok",
			"id":     ticket.ID,
		})
	}
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

//------------------------------------------------------------------------------------------------------------------
// FACEBOOK
//------------------------------------------------------------------------------------------------------------------

func createWorkplaceComment(c *gin.Context, client zendesk.Client) {

}

// Verification Endpoint
func Verify(t string, w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")
	token := r.URL.Query().Get("hub.verify_token")

	if token == os.Getenv(t) {
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Error, wrong validation token"))
	}
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
