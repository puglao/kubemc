package main

import (
	"time"

	"github.com/puglao/kubemc/pkg/kubemc"
)

func main() {
	mcConfig := kubemc.NewKubeMCConfig()

	trigger := make(chan bool, 100)

	go kubemc.WatchMCEvents(mcConfig.KubeMCDir, trigger)
	// fmt.Println(mcConfig)
	// fmt.Println(kubeConfigList)
	for {
		if <-trigger {
			// Avoid unneccessary merge for multi MC file modify event within two seconds
			time.Sleep(mcConfig.MergeRateLimit * time.Second)
			for len(trigger) > 0 {
				<-trigger
			}
			go func() {
				kubemc.MergeKubeMC(mcConfig)
			}()
		}
	}
	<-trigger
}
