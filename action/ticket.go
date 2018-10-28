package ticket

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/MEDIGO/go-zendesk/zendesk"
)

func createAction(jsonData string, client zendesk.Client) {
	rawJsonData, err := ioutil.ReadFile(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	ticketData := &zendesk.Ticket{}
	json.Unmarshal(rawJsonData, &ticketData)

	_, err = client.CreateTicket(ticketData)
	if err != nil {
		log.Fatal(err)
	}
}
