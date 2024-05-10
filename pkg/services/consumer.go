package services

import (
	"main/logs"
	"main/pkg/constants"
	"strings"
)

func StartBussinessLogic() {

	for {
		redisData, ok := <-constants.InternalQueue
		if !ok {
			logs.InfoLog("Found channel closed, BussinessLogic Stopping...")
			constants.ApplicationWaitGroupServices.Done()
			return
		} else {
			logs.InfoLog("Popped redis data from internal channel is: [%v]", redisData)
			if len(redisData) < 2 || len(redisData) > 2 {
				logs.ErrorLog("Invalid redis data received from internal channel: [%v]", redisData)
				//slack alert
			} else {

				qName := redisData[0]
				redisString := redisData[1]
				logs.InfoLog("Valid redis data received from internal channel, queueName: [%v] and  redis string is: [%v]", qName, redisString)

				if strings.EqualFold(strings.ToLower(qName), "MESSAGECHANNEL_MESSAGE_TRANS") {
					PerformmessagechannelMessageTransBussinessLogic(redisString)
				} else if strings.EqualFold(strings.ToLower(qName), "MESSAGECHANNEL_INCOMINGMESSAGE") {
					PerformmessagechannelIncomingBussinessLogic(redisString)
				} else if strings.EqualFold(strings.ToLower(qName), "MESSAGECHANNEL_WEBHOOK_ES") {
					PerformmessagechannelWebhookEsBussinessLogic(redisString)
				} else if strings.EqualFold(strings.ToLower(qName), "UPDATE_CLICK_EVENT_WA") {
					PerformUpdateClickEventWaSingleBussinessLogic(redisString)
				} else if strings.Contains(qName, "MESSAGECHANNEL_OPTIN") {
					PerformmessagechannelOptinBussinessLogic(redisString)
				} else {
					logs.ErrorLog("No matching queue name found for this event:[%v]", redisString)

				}

			}
		}
	}
}
