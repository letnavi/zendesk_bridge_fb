package main

import (
	"encoding/json"
	"log"
	"net/http"

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

const ConfigFile = "conf.toml"

// Processing facebook and zendesk requests
func main() {

	// init config
	var conf Config
	if _, err := toml.DecodeFile(ConfigFile, &conf); err != nil {
		log.Fatal(err)
	}

	// zendesk client
	// sub-domain, email/login and password
	client, err := zendesk.NewClient(conf.Domain, conf.Username, conf.Password)

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

	// verify token fb
	r.GET("/fb", func(c *gin.Context) {
		Verify(conf.Token, c.Writer, c.Request)
	})

	// add comment for facebook
	r.POST("/fb/comment", func(c *gin.Context) {
		createWorkplaceComment(c, client)
	})

	//------------------------------------------------------------------------------------------------------------------
	// ZENDESK
	//------------------------------------------------------------------------------------------------------------------

	// add ticket
	r.POST("/ticket", func(c *gin.Context) {
		if req, err := toZendesk(c); err != nil {
			log.Fatal(err)
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

// request facebook to zendesk decode
func toZendesk(c *gin.Context) ([]byte, error) {
	body, err := c.GetRawData()
	if err != nil {
		log.Fatal(err)
	}
	feed := Feed{}

	json.Unmarshal(body, &feed)

	return feed, nil
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
