package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/logs"
	"main/pkg/constants"
	"net/http"
)

type Slack struct {
	Text string `json:"text"`
}

func BussinessLogicIndex(key string) string {

	logs.InfoLog("indexing:%v", key)

	if len(key) > 0 {
		return string(key[len(key)-1])
	}

	return ""
}

func SendSlackAlert(alert string) {

	todo := Slack{(constants.ApplicationConfig.InstanceName + "\n" + alert)}
	jsonReq, err := json.Marshal(todo)
	resp, err := http.Post("SlackUrl", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}
