package routes

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

// append a message to the messages json  cookie
func SendMessage(c *gin.Context, message string) {
	cookie, _ := c.Cookie("messages")

	// data from json cookie
	var data []string

	err := json.Unmarshal([]byte(cookie), &data)
	if err != nil {
		log.Println("ERROR json unmarshalling for messages failed?", err)
	}

	data = append(data, message)

	jsonData, err1 := json.Marshal(data)
	if err1 != nil {
		log.Println("ERROR json marshalling for messages failed?", err1)
	}

	c.SetCookie("messages", string(jsonData), 5, "/", c.Request.URL.Hostname(), false, false)
}

// get all the json message cookies
func GetMessages(c *gin.Context) []string {
	// get messages json cookie
	messagesCookie, _ := c.Cookie("messages")

	// data from cookie json
	var data []string

	// convert messages json cookie into a map of strings
	err := json.Unmarshal([]byte(messagesCookie), &data)
	if err != nil {
		log.Println("ERROR unmarshalling messages:", err)
	}

	return data
}
