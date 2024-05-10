package services

import (
	"encoding/json"
	"fmt"
	"main/logs"
	"main/pkg/constants"
	database "main/pkg/databases/redis"
	"strings"
	"time"
)

func PerformmessagechannelMessageTransBussinessLogic(redisString string) {

	logs.InfoLog("In PerformmessagechannelMessageTransBussinessLogic for redis string: [%v] ", redisString)
	messagechannelMessageTransMap := make(map[string]interface{})
	json.Unmarshal([]byte(redisString), &messagechannelMessageTransMap)

	for key, value := range messagechannelMessageTransMap {
		logs.InfoLog("Key: %v, Value: %v", key, value)
	}

	// Extract messageJson from the payload
	if messagechannelMessageTransMap["messageJson"] != nil {
		messageJsonMap := make(map[string]interface{})
		// messageJsonStruct := struct {
		// 	Message struct {
		// 		XApiheader string `json:"xApiheader"`
		// 	} `json:"message"`
		// }{}
		messageJsonString := messagechannelMessageTransMap["messageJson"].(string)
		// messageJsonReplaced := strings.Replace(messageJsonString, "\\\\", "\\", -1)
		// fmt.Println("MessageJson String: ", messagechannelMessageTransMap["messageJson"].(string))
		// fmt.Println("MessageJson escapped string:", messageJsonReplaced)
		err := json.Unmarshal([]byte(messageJsonString), &messageJsonMap)
		if err != nil {
			logs.ErrorLog("Error unmarshaling the redis string:[%v], err:[%v]", redisString, err)
			SendSlackAlert(fmt.Sprintf("Error unmarshaling the redis string:[%v], err:[%v]", redisString, err))
		}
		// json.Unmarshal([]byte(messageJsonString), &messageJsonStruct)
		// fmt.Println("MessageJson JSON: ", messageJsonMap)
		// fmt.Println("MessageJson Struct: ", messageJsonStruct)

		for k, v := range messageJsonMap {
			logs.InfoLog("Internal Key %v, Value %v", k, v)
		}

		logs.InfoLog("MessageMapJson: [%v]", messageJsonMap)
		if messagechannelMessageTransMap["errorCode"] == nil && messagechannelMessageTransMap["messagechannelMessageId"] != nil {
			//Set NcMessageId-XapiHeade in Redis

			if messageJsonMap["message"] != nil && messageJsonMap["message"].([]interface{})[0].(map[string]interface{})["xApiheader"] != nil {
				XapiHeader := messageJsonMap["message"].([]interface{})[0].(map[string]interface{})["xApiheader"].(string)
				if XapiHeader != "" {
					logs.InfoLog("XapiHeader [%v]", XapiHeader)
					if messagechannelMessageTransMap["errorCode"] != nil && messagechannelMessageTransMap["messagechannelMessageId"] == nil {
					} else {
						if messagechannelMessageTransMap["messageId"] != nil {
							ncMessageId := messagechannelMessageTransMap["messageId"].(string)
							ncMsgIdXapiHeaderRedisKey := "MESSAGECHANNEL_X_API_HEADERS_" + ncMessageId
							database.CustomSetKey(ncMsgIdXapiHeaderRedisKey, XapiHeader, time.Hour*24*5)
						} else {
							logs.InfoLog("Nc message Id not found for this payload [%v]", redisString)
						}
					}
				}

			}

			if messagechannelMessageTransMap["CLIENT"] != nil {
				clientId_phoneNumber := messagechannelMessageTransMap["CLIENT"].(string)
				clientId := strings.Split(clientId_phoneNumber, "_")[0]
				logs.InfoLog("Client ID: %v", clientId)

				currentTime := time.Now()
				formattedDate := currentTime.Format("02012006")
				hashKey := "MESSAGECHANNEL_MESSAGE_ID_MAPPING_" + clientId + "_" + formattedDate
				hashKey1 := "NC_REQUEST_MESSAGE_ID_MAPPING_" + clientId + "_" + formattedDate
				messageId := messagechannelMessageTransMap["messagechannelMessageId"].(string)
				NcMessageId := messagechannelMessageTransMap["messageId"].(string)

				//HSet MESSAGECHANNEL_MESSAGE_ID_MAPPING_<clientId>_<date> messageId NcMsgID_<date>
				database.CustomHSet(hashKey, messageId, NcMessageId)

				//HSet NC_REQUEST_MESSAGE_ID_MAPPING_<clientId>_<date>
				database.CustomHSet(hashKey1, NcMessageId, messageId)
			}

		} else {
			logs.InfoLog("No redis operation for this payload [%v]", redisString)
		}

		//Extract the phone number from the payload, phoneNumber  will provide the BussinessLogic index
		listName := "MESSAGECHANNEL_MESSAGE"
		index := "0"
		if messageJsonMap["phoneNumber"] != nil {
			phoneNumberString := messageJsonMap["phoneNumber"].(string)
			index = BussinessLogicIndex(phoneNumberString)
			logs.InfoLog("Index for payload [%v], having phone number [%v] is [%v]", redisString, phoneNumberString, index)

			if (messagechannelMessageTransMap["errorCode"] != nil && messagechannelMessageTransMap["errorCode"] != "") || (messagechannelMessageTransMap["errorMessage"] != nil && messagechannelMessageTransMap["errorMessage"] != "") {
				listName = listName + "_FAILED" + "_" + index
				database.CustomRpush(listName, redisString)
				if constants.GcpBussinessLogic {
					database.CustomRpush(listName+"_GCP", redisString)
				}
			} else {
				listName = listName + "_" + "PUBLISH_" + index
				database.CustomRpush(listName, redisString)
				if constants.GcpBussinessLogic {
					database.CustomRpush(listName+"_GCP", redisString)
				}
			}
		} else {
			logs.ErrorLog("Error pushing this redis string to segregated queue, couldnt determine the index [%v]", redisString)
			database.CustomRpush("MESSAGECHANNEL_MESSAGE_PARKED", redisString)
		}

	} else {
		logs.InfoLog("No message to extract from payload [%v]", redisString)
	}
}

func PerformmessagechannelWebhookEsBussinessLogic(redisString string) {
	messagechannelWebhookEsMap := make(map[string]interface{})
	json.Unmarshal([]byte(redisString), &messagechannelWebhookEsMap)
	logs.InfoLog("Map for redis string [%v], messagechannelWebhookEsMap [%v]", redisString, messagechannelWebhookEsMap)
	if messagechannelWebhookEsMap["phoneNumber"] != nil && messagechannelWebhookEsMap["webhookJson"] != nil {

		logs.InfoLog("Phone number extracted from payload [%v] is [%v]", messagechannelWebhookEsMap, messagechannelWebhookEsMap["phoneNumber"].(string))

		logs.InfoLog("WebhookJson extracted from payload [%v] is [%v]", messagechannelWebhookEsMap, messagechannelWebhookEsMap["webhookJson"].(string))

		webhookJsonMap := make(map[string]interface{})
		json.Unmarshal([]byte(messagechannelWebhookEsMap["webhookJson"].(string)), &webhookJsonMap)

		index := BussinessLogicIndex(messagechannelWebhookEsMap["phoneNumber"].(string))

		listName := "MESSAGECHANNEL_WEBHOOK_ES"
		for _, status := range constants.ApplicationConfig.Application.WebhookEventToSegregate {
			if status == webhookJsonMap["statuses"].([]interface{})[0].(map[string]interface{})["status"].(string) {
				listName = listName + "_" + status
			}
		}

		listName = listName + "_" + index
		database.CustomRpush(listName, redisString)

		if constants.GcpBussinessLogic {
			database.CustomRpush(listName+"_GCP", redisString)
		}
	}

}

func PerformmessagechannelIncomingBussinessLogic(redisString string) {
	messagechannelIncomingEsMap := make(map[string]interface{})
	json.Unmarshal([]byte(redisString), &messagechannelIncomingEsMap)
	logs.InfoLog("Map for redis string [%v], messagechannelIncomingMessageMap [%v]", redisString, messagechannelIncomingEsMap)
	index := BussinessLogicIndex(messagechannelIncomingEsMap["phoneNumber"].(string))
	listName := "MESSAGECHANNEL_INCOMINGMESSAGE"
	listName = listName + "_" + index
	database.CustomRpush(listName, redisString)
	if constants.GcpBussinessLogic {
		database.CustomRpush(listName+"_GCP", redisString)
	}
}

func PerformUpdateClickEventWaSingleBussinessLogic(redisString string) {
	listName := "UPDATE_CLICK_EVENT_WA"
	database.CustomRpush(listName+"_AWS", redisString)

	if constants.GcpBussinessLogic {
		database.CustomRpush(listName+"_GCP", redisString)
	}

}

func PerformmessagechannelOptinBussinessLogic(redisString string) {
	listName := "MESSAGECHANNEL_OPTIN"

	database.CustomRpush(listName+"_AWS", redisString)

	if constants.GcpBussinessLogic {
		database.CustomRpush(listName+"_GCP", redisString)
	}
}
