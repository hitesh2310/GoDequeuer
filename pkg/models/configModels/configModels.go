package models

type Config struct {
	InstanceName string `json:"instanceName"`
	Database     struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
	} `json:"database"`

	RedisQueue struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
	} `json:"redisQueue"`

	RedisCache struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
	} `json:"redisCache"`

	Application struct {
		LogPath                 string   `json:"logPath"`
		Chanlen                 string   `json:"chanlen"`
		Pullers                 string   `json:"pullers"`
		Workers                 string   `json:"workers"`
		MasterCacheKey          string   `json:"master_cache_key"`
		WaUsrMap                string   `json:"wa_usr_map"`
		WaPhNoMap               string   `json:"wa_ph_no_map"`
		GcpBussinessLogic       string   `json:"gcpBussinessLogic"`
		WebhookEventToSegregate []string `json:"webhookEventToSegregate"`
	} `json:"application"`
}
