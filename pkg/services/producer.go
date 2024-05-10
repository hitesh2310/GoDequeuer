package services

import (
	"main/logs"
	"main/pkg/constants"
	database "main/pkg/databases/redis"
)

func StartPuller(qname string) {
	for {
		if constants.Shutdown {
			logs.InfoLog("Puller of queue :[%v] is shutting down....", qname)
			logs.InfoLog("Closing channel... ")
			//if not closed already:
			if !constants.Closed {
				close(constants.InternalQueue)
				constants.Closed = true
			}

			constants.ApplicationWaitGroupServices.Done()
			return
		} else {
			// logs.InfoLog("Puller of queue :[%v]", qname)
			result, err := database.CustomBLpop(qname)

			if err != nil && err.Error() == "redis: nil" {
				//do nothing
			} else if err != nil {
				logs.ErrorLog("Error to pull redis string from queue: [%v], got error:[%v] ", qname, err)
				//send alert
			} else {
				logs.InfoLog("Pulled redis string from queue: [%v] is: [%v]", result[0], result[1])

				if len(constants.InternalQueue) < 1 {
					constants.InternalQueue <- result
				} else {
					logs.InfoLog("InternalQueue is full repusing the redis string data  %v", result)
					database.CustomRpush(result[0], result[1])
				}

			}
		}
	}
}
