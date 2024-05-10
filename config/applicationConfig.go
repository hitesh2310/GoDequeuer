package config

import (
	"encoding/json"
	"fmt"
	"main/logs"
	"main/pkg/constants"
	database "main/pkg/databases/redis"
	models "main/pkg/models/configModels"
	"os"
	"strconv"
)

func SetUpApplication() {
	fmt.Println("SetUp config details....")
	setupConfig()
	fmt.Println("SettingUp Logs....")
	setUpApplicationLogs()
	setUpApplicationDatabase()
	fmt.Println("DB set up done")
	Chansize, err := strconv.Atoi(constants.ApplicationConfig.Application.Chanlen)
	if err != nil {
		constants.Chansize = 1
	} else {
		constants.Chansize = Chansize
	}

	constants.InternalQueue = make(chan []string, constants.Chansize)

	NumberOfPullers, err := strconv.Atoi(constants.ApplicationConfig.Application.Pullers)
	if err != nil {
		constants.Pullers = 1
	} else {
		constants.Pullers = NumberOfPullers
	}

	NumberOfWorkers, err := strconv.Atoi(constants.ApplicationConfig.Application.Pullers)
	if err != nil {
		constants.Workers = 2
	} else {
		constants.Workers = NumberOfWorkers
	}

	// //GLOBAL CLIENT-ID SET
	// clientMap := make(map[string]bool)
	// constants.CampaignInProcess = clientMap

	//BATCH SIZE
	// batchSizeStr := constants.ApplicationConfig.Application.BatchSize
	// batchSizeInt, err := strconv.Atoi(batchSizeStr)
	// if err != nil {
	// 	batchSizeInt = 10000
	// }
	// constants.BatchSize = batchSizeInt
	// fmt.Println("API max BatchSize: ", constants.BatchSize)

	// //offStart and offEnd  time
	// fmt.Println("OFF START :", constants.ApplicationConfig.Application.OffTimeStart)
	// fmt.Println("OFF END :", constants.ApplicationConfig.Application.OffTimeEnd)
}

func setupConfig() {
	// viper.SetConfigName("config")    // Name of the config file (without extension)
	// viper.SetConfigType("json")      // Config file type
	// viper.AddConfigPath("./config/") // Path to the directory containing the config file

	// // Read the config file
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	panic(fmt.Errorf("failed to read config file: %s", err))
	// }

	// var config models.Config
	// err = viper.Unmarshal(&config)
	// if err != nil {
	// 	fmt.Println("Error to unmarshal config")
	// }

	file, err := os.ReadFile("./config/config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Unmarshal the JSON data into a Config struct
	var config models.Config

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	constants.ApplicationConfig = &config
	fmt.Println("config.Application.GcpBussinessLogic", config.Application.GcpBussinessLogic)

	constants.GcpBussinessLogic, _ = strconv.ParseBool(config.Application.GcpBussinessLogic)

	fmt.Println("GCP value:", constants.GcpBussinessLogic)

	fmt.Println("Events:", config.Application.WebhookEventToSegregate)
}

func setUpApplicationLogs() {
	logs.SetUpApplicationLogs()
}

func setUpApplicationDatabase() {
	// databaseMySql.EstablishDbConnection()
	database.EstablishRedisCacheConnecion()
	database.EstablishRedisQueueConnecion()
}
