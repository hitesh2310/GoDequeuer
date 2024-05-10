package constants

import (
	models "main/pkg/models/configModels"

	"sync"
)

var (
	ApplicationConfig            *models.Config
	ApplicationWaitGroupServices = sync.WaitGroup{}
	InternalQueue                chan []string
	Shutdown                     bool
	BatchSize                    int
	Chansize                     int
	Pullers                      int
	Workers                      int
	DoneSignal                   chan struct{}
	GcpBussinessLogic            bool
	Closed                       bool
)
