package services

import (
	"fmt"
	"main/pkg/constants"

	"os"
	"os/signal"
	"syscall"
)

func StartServices(qList []string) {

	fmt.Println("Starting services ..... ")
	for i := 0; i < constants.Pullers; i++ {
		// Starting PUllers
		for _, queue := range qList {
			constants.ApplicationWaitGroupServices.Add(1)
			go StartPuller(queue)
		}

	}

	for i := 0; i < constants.Workers; i++ {
		constants.ApplicationWaitGroupServices.Add(1)
		//start BussinessLogics
		go StartBussinessLogic()
	}
	serviceShutDown := make(chan os.Signal, 1)
	signal.Notify(serviceShutDown, os.Interrupt, syscall.SIGTERM)

	<-serviceShutDown
	fmt.Println("Stopping services .....")
	constants.Shutdown = true

	fmt.Println("Waiting for services.. ")
	constants.ApplicationWaitGroupServices.Wait()
}
