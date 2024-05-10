package main

import (
	"fmt"
	"main/config"
	"main/pkg/services"
	"os"
)

func init() {
	fmt.Println("Starting messagechannelBussinessLogic....")
	config.SetUpApplication()

}

func main() {
	parameters := os.Args[1:]
	fmt.Println("in main startingservices...")
	fmt.Println("parameters:", parameters)
	services.StartServices(parameters)
}
